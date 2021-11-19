package clientctl

import (
	"context"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/zxfonline/IMDemo/config"
	"github.com/zxfonline/IMDemo/core/session"
	"github.com/zxfonline/IMDemo/model"
)

type ClientServer struct {
	// 下线、掉线的玩家
	LogoutChan chan int64
	Rooms      map[int64]*model.ChatRoom
}

var (
	SvrCtl = &ClientServer{
		LogoutChan: make(chan int64, 0x1000),
		Rooms:      make(map[int64]*model.ChatRoom, 4),
	}
)

func StartServer(ctx context.Context, wg *sync.WaitGroup, roomSize int64, chatCashSize int32) {
	SvrCtl.Start(ctx, wg, roomSize, chatCashSize)
}

// start server loop
func (s *ClientServer) Start(ctx context.Context, wg *sync.WaitGroup, roomSize int64, chatCashSize int32) {
	for i := int64(1); i <= roomSize; i++ {
		room := model.NewChatRoom(i, chatCashSize)
		go room.Run(ctx, wg)
		SvrCtl.Rooms[i] = room

	}
	go s.handleMsg(ctx, wg)
}

func (s *ClientServer) ClientLogic(ctxt context.Context, wg *sync.WaitGroup, conn *websocket.Conn) {
	// 创建会话
	msgChan := make(chan *session.NetPacket, 30)
	sendChan := make(chan *session.NetPacket, 256)
	session := session.NewSession(conn, msgChan, sendChan, s.LogoutChan)
	session.SetParameter(30*time.Second, 30*time.Second, 512, true)

	if config.IsDebug() {
		session.SetRpmParameter(0, 0, nil)
	} else {
		session.SetRpmParameter(36, 3*time.Second, nil)
	}

	agent := model.NewClientAgent(session)
	model.ClientAgentAdd(agent)
	session.HandleConn(nil)
	go handleServerMsg(ctxt, wg, agent)
}

// 消息分发
func (s *ClientServer) handleMsg(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case sessionId := <-s.LogoutChan: // 连接掉线
			if userAgent := model.ClientAgentGet(sessionId); userAgent != nil {
				if roomID := userAgent.State.Load(); roomID > 0 { //当前在房间
					room := s.Rooms[roomID]
					room.Unregister <- userAgent
				}
				model.ClientAgentOffline(sessionId)
			}
		}
	}
}

//随机获取一个聊天房间
func (s *ClientServer) RandRoom() *model.ChatRoom {
	for _, room := range s.Rooms {
		return room
	}
	return nil
}

//获取聊天房间
func (s *ClientServer) Room(roomID int64) *model.ChatRoom {
	return s.Rooms[roomID]
}

// 消息分发
func handleServerMsg(ctxt context.Context, wg *sync.WaitGroup, client *model.ClientAgent) {
	wg.Add(1)
	defer wg.Done()
	for {
		select {
		case <-ctxt.Done():
			return
		case <-client.Session.CloseState:
			return
		case msg := <-client.Session.ReadChan:
			handleMsg(ctxt, wg, client, msg)
		}
	}
}
