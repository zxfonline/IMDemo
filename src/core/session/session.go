package session

import (
	"bytes"
	"errors"
	"net"
	"sync/atomic"
	"time"

	"github.com/zxfonline/IMDemo/core/chanutil"

	"github.com/gorilla/websocket"
	"github.com/zxfonline/IMDemo/core/log"
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)
var _sessionID int64

var WSUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type NetPacket struct {
	MsgType interface{}
	Data    []byte

	//收到该消息包的时间戳 毫秒
	ReceiveTime time.Time
}

func NewSession(conn *websocket.Conn, readChan, sendChan chan *NetPacket, offChan chan int64) *WsSession {
	now := time.Now()
	s := &WsSession{
		Conn:          conn,
		SendChan:      sendChan,
		ReadChan:      readChan,
		OffChan:       offChan,
		SessionId:     atomic.AddInt64(&_sessionID, 1),
		sendDelay:     10 * time.Second,
		readDelay:     60 * time.Second,
		pingPeriod:    60 * 9 / 10 * time.Second,
		sendFullClose: true,
		OnLineTime:    &now,
		CloseState:    chanutil.NewDoneChan(),
	}
	s.Conn.SetReadLimit(1024)
	// log.Debugf("new connection from:%v", conn.RemoteAddr().String())
	return s
}

type WsSession struct {
	Conn     *websocket.Conn
	SendChan chan *NetPacket
	ReadChan chan *NetPacket
	//离线消息管道,用于外部接收连接断开的消息并处理后续
	OffChan chan int64
	// ID
	SessionId int64

	readDelay  time.Duration
	sendDelay  time.Duration
	pingPeriod time.Duration
	// Declares how many times we will try to resend message
	MaxSendRetries int
	//发送管道满后是否需要关闭连接
	sendFullClose bool
	CloseState    chanutil.DoneChan

	// 包频率包数
	rpmLimit uint32
	// 包频率检测间隔
	rpmInterval time.Duration
	// 超过频率控制离线通知包
	rpmLimitMsg *NetPacket

	//登录时间
	OnLineTime *time.Time
	//离线时间
	OffLineTime *time.Time
}

//filter:true 过滤成功，抛弃该报文；false:过滤失败，继续执行该报文消息
func (s *WsSession) HandleConn(filter func(*NetPacket) bool) {
	go s.ReadLoop(filter)
	go s.SendLoop()
}

//网络连接远程ip
func (s *WsSession) RemoteAddr() net.Addr {
	return s.Conn.RemoteAddr()
}

func (s *WsSession) Send(packet *NetPacket) bool {
	if packet == nil {
		return false
	}
	if !s.sendFullClose { //阻塞发送，直到管道关闭
		select {
		case s.SendChan <- packet:
			if wait := len(s.SendChan); wait > cap(s.SendChan)/10*5 && wait%20 == 0 {
				log.Warnf("session send process,waitChan:%d/%d,msg:%v,session:%d,remote:%s", wait, cap(s.SendChan), packet.MsgType, s.SessionId, s.RemoteAddr())
			}
			return true
		case <-s.CloseState:
			return false
		}
	} else { //缓存管道满了会关闭连接
		select {
		case <-s.CloseState:
			return false
		case s.SendChan <- packet:
			if wait := len(s.SendChan); wait > cap(s.SendChan)/10*5 && wait%20 == 0 {
				log.Warnf("session send process,waitChan:%d/%d,msg:%v,session:%d,remote:%s", wait, cap(s.SendChan), packet.MsgType, s.SessionId, s.RemoteAddr())
			}
			return true
		default:
			log.Errorf("session sender overflow,close session,waitChan:%d,msg:%v,session:%d,remote:%s", len(s.SendChan), packet.MsgType, s.SessionId, s.RemoteAddr())
			s.Close()
			return false
		}
	}
}

func (s *WsSession) ReadLoop(filter func(*NetPacket) bool) {
	defer log.PrintPanicStack()

	// 关闭发送
	defer s.Close()

	rpmStart := time.Now()
	rpmCount := uint32(0)
	s.Conn.SetPongHandler(func(string) error {
		if s.readDelay > 0 {
			s.Conn.SetReadDeadline(time.Now().Add(s.readDelay))
		}
		return nil
	})

	//rpmMsgCount := 0

	for {
		// 读取超时
		if s.readDelay > 0 {
			s.Conn.SetReadDeadline(time.Now().Add(s.readDelay))
		}

		messageType, message, err := s.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Error(err)
			}
			return
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		// 收包频率控制
		if s.rpmLimit > 0 {
			rpmCount++

			// 达到限制包数
			if rpmCount > s.rpmLimit {
				now := time.Now()
				// 检测时间间隔
				if now.Sub(rpmStart) < s.rpmInterval {
					// 提示操作太频繁三次后踢下线
					//rpmMsgCount++
					//if rpmMsgCount > 3 {
					s.DirectSendAndClose(s.rpmLimitMsg)
					log.Errorf("session rpm too high,%d/%s qps,session:%d,remote:%s", rpmCount, s.rpmInterval, s.SessionId, s.RemoteAddr())
					return
					//}
				}

				// 未超过限制
				rpmCount = 0
				rpmStart = now
			}
		}

		pack := &NetPacket{MsgType: messageType, Data: message, ReceiveTime: time.Now()}

		if filter == nil {
			s.ReadChan <- pack
		} else {
			if ok := filter(pack); !ok {
				s.ReadChan <- pack
			}
		}
	}
}

func (s *WsSession) Close() {
	s.CloseState.SetDone()
}

func (s *WsSession) IsClosed() bool {
	return s.CloseState.R().Done()
}

func (s *WsSession) closeTask() {
	offTime := time.Now()
	s.OffLineTime = &offTime
	if s.OffChan != nil {
		s.OffChan <- s.SessionId
	}
	s.Conn.Close()
}

func (s *WsSession) SendLoop() {
	defer log.PrintPanicStack()
	pingPeriod := 1 * time.Minute
	if s.pingPeriod > 0 {
		pingPeriod = s.pingPeriod
	}
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
	}()
	for {
		select {
		case <-s.CloseState:
			s.closeTask()
			return
		case <-ticker.C:
			s.DirectSend(&NetPacket{MsgType: websocket.PingMessage})
		case packet := <-s.SendChan:
			s.DirectSend(packet)
		}
	}
}

func (s *WsSession) DirectSendAndClose(packet *NetPacket) {
	go func() {
		if s.DirectSend(packet) {
			time.Sleep(1 * time.Second)
			s.Close()
		}
	}()
}

func (s *WsSession) DirectSend(packet *NetPacket) bool {
	if packet == nil {
		return true
	}
	if s.IsClosed() {
		return false
	}
	err := s.performSend(packet, 0)
	if err != nil {
		log.Debugf("error writing msg,session:%d,remote:%s,err:%v", s.SessionId, s.RemoteAddr(), err)
		s.Close()
		return false
	}
	return true
}

func (s *WsSession) performSend(msg *NetPacket, sendRetries int) error {
	// 写超时
	if s.sendDelay > 0 {
		s.Conn.SetWriteDeadline(time.Now().Add(s.sendDelay))
	}
	switch msg.MsgType {
	case websocket.TextMessage, websocket.BinaryMessage:
	case websocket.CloseMessage, websocket.PingMessage, websocket.PongMessage:
	default:
		return errors.New("unsupport message type")
	}
	err := s.Conn.WriteMessage(msg.MsgType.(int), msg.Data)
	if err != nil {
		return s.processSendError(err, msg, sendRetries)
	}
	return nil
}

func (s *WsSession) processSendError(err error, msg *NetPacket, sendRetries int) error {
	netErr, ok := err.(net.Error)
	if !ok {
		return err
	}

	if s.isNeedToResendMessage(netErr, sendRetries) {
		return s.performSend(msg, sendRetries+1)
	}
	return err
}

func (s *WsSession) isNeedToResendMessage(err net.Error, sendRetries int) bool {
	return (err.Temporary() || err.Timeout()) && sendRetries < s.MaxSendRetries
}

// 设置链接参数
func (s *WsSession) SetParameter(readDelay, sendDelay time.Duration, maxRecvSize uint32, sendFullClose bool) {
	s.Conn.SetReadLimit(int64(maxRecvSize))
	if readDelay >= 0 {
		s.readDelay = readDelay
		s.pingPeriod = readDelay * 9 / 10
	}
	if sendDelay >= 0 {
		s.sendDelay = sendDelay
	}
	s.sendFullClose = sendFullClose
}

// 包频率控制参数
func (s *WsSession) SetRpmParameter(rpmLimit uint32, rpmInterval time.Duration, msg *NetPacket) {
	s.rpmLimit = rpmLimit
	s.rpmInterval = rpmInterval
	s.rpmLimitMsg = msg
	if s.rpmLimitMsg == nil {
		s.rpmLimitMsg = &NetPacket{
			MsgType: websocket.CloseMessage,
			Data:    websocket.FormatCloseMessage(websocket.CloseNormalClosure, "messages are sent too frequently"),
		}
	}
}
