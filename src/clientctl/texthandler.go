package clientctl

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/valyala/fastjson"
	"github.com/zxfonline/IMDemo/core/badword"
	"github.com/zxfonline/IMDemo/core/gerror"
	"github.com/zxfonline/IMDemo/core/log"
	"github.com/zxfonline/IMDemo/core/nametrie"
	"github.com/zxfonline/IMDemo/core/session"
	"github.com/zxfonline/IMDemo/model"
)

var (
	NameReapCheck = nametrie.NewNameMatchTrie()
)

func ProcessTextMessage(ctx context.Context, wg *sync.WaitGroup, clientAgent *model.ClientAgent, msg *session.NetPacket) (err error, retMsg []*session.NetPacket) {
	v, perr := fastjson.ParseBytes(msg.Data)
	if perr != nil {
		err = perr
		retMsg = []*session.NetPacket{{
			MsgType: websocket.CloseMessage,
			Data:    websocket.FormatCloseMessage(websocket.CloseInvalidFramePayloadData, "parse payload err"),
		}}
		return
	}
	log.Debugf("request:%v", v.String())
	//Request
	reqType := v.GetUint("type")
	ackType := reqType + 1
	errCode := gerror.SERVER_CDATA_ERROR
	//处理异常错误
	defer func() {
		if err != nil {
			retMsg = []*session.NetPacket{{
				MsgType: websocket.TextMessage,
				Data: (&Response{
					Type:    RequestType(ackType),
					Code:    errCode,
					CodeMsg: err.Error(),
				}).toJson(),
			}}
		}
	}()
	//捕获异常
	defer gerror.PanicToErr(&err)

	roomIDState := clientAgent.State.Load()
	switch RequestType(reqType) {
	case LoginReq:
		ackType = uint(LoginAck)
		if roomIDState == -1 { //玩家掉线
			err = errors.New("you are logout,refresh page(F5)")
			return
		}
		if roomIDState > 0 { //玩家已经有房间了
			errCode = ERROR_IGNORE //不做操作
			err = errors.New("you are in the chat room")
			return
		}
		userName := v.GetStringBytes("data", "userName")
		if NameReapCheck.FullMatch(string(userName)) {
			errCode = ERROR_NAME_REPEAT //姓名重复
			err = errors.New("repeated name,change name please")
			return
		}
		NameReapCheck.Add(string(userName))
		clientAgent.UserName = string(userName)
		room := SvrCtl.RandRoom()
		clientAgent.State.Store(room.RoomID)

		retMsg = []*session.NetPacket{{
			MsgType: websocket.TextMessage,
			Data: (&Response{
				Type: RequestType(ackType),
				Code: gerror.OK,
				Data: &struct {
					RoomID   int64  `json:"roomID"`
					UserName string `json:"userName"`
				}{
					RoomID:   room.RoomID,
					UserName: clientAgent.UserName,
				},
			}).toJson(),
		}}
		retMsg = append(retMsg, room.RecentMsg...)
		room.Register <- clientAgent
		return
	case RoomSwitchReq:
		ackType = uint(RoomSwitchAck)
		oldRoom := SvrCtl.Room(roomIDState)
		if oldRoom == nil {
			err = errors.New("you haven not logged in yet")
			return
		}
		roomID := v.GetInt64("data", "room")
		newRoom := SvrCtl.Room(roomID)
		if newRoom == nil {
			err = errors.New("not found new room")
			return
		}
		if roomID == roomIDState { //房间相同不用处理
			return
		}
		//更换房间
		clientAgent.State.Store(newRoom.RoomID)
		oldRoom.Unregister <- clientAgent

		retMsg = []*session.NetPacket{{
			MsgType: websocket.TextMessage,
			Data: (&Response{
				Type: RequestType(ackType),
				Code: gerror.OK,
				Data: &struct {
					RoomID   int64  `json:"roomID"`
					UserName string `json:"userName"`
				}{
					RoomID:   newRoom.RoomID,
					UserName: clientAgent.UserName,
				},
			}).toJson(),
		}}
		retMsg = append(retMsg, newRoom.RecentMsg...)
		newRoom.Register <- clientAgent
		return
	case RoomChatReq: // 当前房间聊天
		ackType = uint(RoomChatAck)
		if roomIDState > 0 { //当前房间
			//构建消息用户名和发送时间,替换脏字
			chatMessage := string(v.GetStringBytes("data", "message"))
			chatMessage = badword.BadWordReplace(chatMessage)
			v.Get("data").Set("message", fastjson.MustParse(fmt.Sprintf("%q", chatMessage)))
			v.Get("data").Set("userName", fastjson.MustParse(fmt.Sprintf("%q", clientAgent.UserName)))
			v.Get("data").Set("sendTime", fastjson.MustParse(fmt.Sprintf("%q", time.Now().Format("2006-01-02 15:04:05"))))
			v.Set("type", fastjson.MustParse(fmt.Sprintf("%d", RoomChatNtf)))
			msg.Data = []byte(v.String())
			room := SvrCtl.Room(roomIDState)
			room.Broadcast <- msg
			return
		} else {
			err = errors.New("no found chat room,refresh page(F5)")
			return
		}
	}
	return
}
