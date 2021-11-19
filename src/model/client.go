package model

import (
	"sync"

	"github.com/zxfonline/IMDemo/core/atomic"
	"github.com/zxfonline/IMDemo/core/log"
	"github.com/zxfonline/IMDemo/core/session"
)

var (
	//key=sessionID value=*ClientAgent
	_clientKv sync.Map
)

type ClientAgent struct {
	Session  *session.WsSession
	UserName string
	//-1掉线,0大厅,1,2,3...房间id
	State *atomic.Int64
}

func NewClientAgent(session *session.WsSession) *ClientAgent {
	return &ClientAgent{
		Session: session,
		State:   atomic.NewInt64(0),
	}
}

func ClientAgentAdd(client *ClientAgent) {
	_clientKv.Store(client.Session.SessionId, client)
	log.Debugf("connected client,session:%d,remote:%s", client.Session.SessionId, client.Session.RemoteAddr())
}

func ClientAgentGet(sessionID int64) *ClientAgent {
	if tmp, ok := _clientKv.Load(sessionID); ok {
		client := tmp.(*ClientAgent)
		return client
	}
	return nil
}

func RangeSessions(callback func(clientAgent *ClientAgent) bool) {
	_clientKv.Range(func(key, session interface{}) bool {
		cs := session.(*ClientAgent)
		if cs.Session.IsClosed() || cs.State.Load() == -1 {
			return true
		}
		return callback(cs)
	})
}

func ClientAgentOffline(sessionID int64) int64 {
	defer log.PrintPanicStack()
	if tmp, ok := _clientKv.Load(sessionID); ok {
		client := tmp.(*ClientAgent)
		_clientKv.Delete(sessionID)
		client.State.Store(-1)
	}
	return 0
}
