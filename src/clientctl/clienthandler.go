package clientctl

import (
	"context"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/zxfonline/IMDemo/core/log"
	"github.com/zxfonline/IMDemo/core/session"
	"github.com/zxfonline/IMDemo/model"
)

func handleMsg(ctx context.Context, wg *sync.WaitGroup, clientAgent *model.ClientAgent, msg *session.NetPacket) {
	defer log.PrintPanicStack()
	select {
	case <-ctx.Done():
		return
	case <-clientAgent.Session.CloseState:
		return
	default:
		if msg.MsgType == websocket.TextMessage {
			err, retMsgs := ProcessTextMessage(ctx, wg, clientAgent, msg)
			if len(retMsgs) != 0 {
				for _, retMsg := range retMsgs {
					log.Debugf("send:%v", string(retMsg.Data))
					clientAgent.Session.Send(retMsg)
				}
			}
			if err != nil {
				log.Errorf("process msg err:%v", err)
			}
		} else {
			log.Warnf("unsupport message type:%v", msg.MsgType)
			clientAgent.Session.DirectSend(&session.NetPacket{
				MsgType: websocket.CloseMessage,
				Data:    websocket.FormatCloseMessage(websocket.CloseUnsupportedData, ""),
			})
		}

	}
}
