package model

import (
	"context"
	"sync"
	"time"

	"github.com/valyala/fastjson"
	"github.com/zxfonline/IMDemo/core/hotword"
	"github.com/zxfonline/IMDemo/core/log"
	"github.com/zxfonline/IMDemo/core/session"
)

type ChatRoom struct {
	RoomID     int64
	clients    map[*ClientAgent]bool
	Broadcast  chan *session.NetPacket
	Register   chan *ClientAgent
	Unregister chan *ClientAgent
	//最新的缓存消息
	RecentMsg []*session.NetPacket
	//热门消息记录
	HotMsg *hotword.TimeTrie
}

func NewChatRoom(roomID int64, cacheChatSize int32) *ChatRoom {
	room := &ChatRoom{
		RoomID:     roomID,
		clients:    make(map[*ClientAgent]bool, 128),
		Broadcast:  make(chan *session.NetPacket, 1024),
		Register:   make(chan *ClientAgent, 16),
		Unregister: make(chan *ClientAgent, 16),
		RecentMsg:  make([]*session.NetPacket, 0, cacheChatSize),
		HotMsg:     hotword.NewTimeTrie(),
	}
	return room
}

func (cr *ChatRoom) Run(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	var interval int64 = 60
	realExpire := interval - (time.Now().Unix() % interval)
	ticker := time.NewTimer(time.Duration(realExpire) * time.Second)
	defer func() {
		wg.Done()
		ticker.Stop()
	}()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			ticker.Reset(time.Duration(interval) * time.Second)
			//清理当前时间十分钟以前的热词信息
			cr.HotMsg.OnTimeout(time.Now().Unix() - 10*60)
		case client := <-cr.Register:
			cr.clients[client] = true
		case client := <-cr.Unregister:
			delete(cr.clients, client)
		case message := <-cr.Broadcast:
			cr.broadcastLogic(message)
		}
	}
}
func (cr *ChatRoom) broadcastLogic(message *session.NetPacket) {
	defer log.PrintPanicStack()
	cr.addRecentMsg(message)
	for client := range cr.clients {
		if client.State.Load() == cr.RoomID {
			client.Session.Send(message)
		}
	}
}
func (cr *ChatRoom) addRecentMsg(msg *session.NetPacket) {
	arr := cr.RecentMsg
	if len(arr) >= cap(arr) {
		copy(arr[:], arr[1:])
		arr[cap(arr)-1] = msg
	} else {
		arr = append(arr, msg)
	}
	cr.RecentMsg = arr
	if chat := fastjson.GetString(msg.Data, "data", "message"); chat != "" {
		//FIXME 可以考虑使用分词器进行优化，目前默认将一句话作为热词
		cr.HotMsg.Add(chat)
	}

}
