package service

import (
	"context"
	"sync"
	"time"

	"github.com/zxfonline/IMDemo/clientctl"
	"github.com/zxfonline/IMDemo/core/gerror"
	"github.com/zxfonline/IMDemo/core/log"
	"github.com/zxfonline/IMDemo/core/strutil"
	"github.com/zxfonline/IMDemo/model"

	"github.com/zxfonline/IMDemo/core/session"
	"github.com/zxfonline/IMDemo/core/web"
)

func RegisterHandlers(ctxt context.Context, wg *sync.WaitGroup, server *web.Server) {
	server.Get("/chat", func(ctx *web.Context) (interface{}, error) {
		conn, err := session.WSUpgrader.Upgrade(ctx.ResponseWriter, ctx.Request, nil)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		clientctl.SvrCtl.ClientLogic(ctxt, wg, conn)
		return nil, nil
	})
	//当前最热的话 `/popular/(房间号1-4)`
	server.Get("/popular/([1-9]+)/([1-9]\\d*)", func(ctx *web.Context, room string, topX string) (interface{}, error) {
		roomID := strutil.Stoi64(room, 0)
		hotNum := strutil.Stoi(topX, 1)
		roomInfo := clientctl.SvrCtl.Room(roomID)
		if roomInfo == nil {
			return nil, gerror.NewError(gerror.SERVER_CMSG_ERROR, "no room found")
		}
		hots := roomInfo.HotMsg.HotTopX(hotNum)
		return &struct {
			Code int         `json:"code"`
			Data interface{} `json:"data,omitempty"`
		}{
			Code: int(gerror.OK),
			Data: hots,
		}, nil
	})
	//查询在线玩家的信息 `/stats/(角色名)`
	server.Get("/stats", func(ctx *web.Context) (interface{}, error) {
		name := ctx.Param("name", "")
		var clientAgent *model.ClientAgent
		if name != "" {
			model.RangeSessions(func(findAgent *model.ClientAgent) bool {
				if findAgent.UserName == name {
					clientAgent = findAgent
					return false
				}
				return true
			})
		}
		if clientAgent == nil {
			return nil, gerror.NewError(gerror.SERVER_CMSG_ERROR, "no player found")
		}
		now := time.Now()
		lt := clientAgent.Session.OnLineTime
		if lt == nil {
			lt = &now
		}
		loginTime := lt.Format("2006-01-02 15:04:05")
		return &struct {
			Code       int    `json:"code"`
			UserName   string `json:"userName"`
			LoginTime  string `json:"loginTime"`
			OnlineTime string `json:"onlineTime"`
			RoomID     int64  `json:"roomID"`
		}{
			Code:       int(gerror.OK),
			UserName:   clientAgent.UserName,
			LoginTime:  loginTime,
			OnlineTime: time.Since(*lt).String(),
			RoomID:     clientAgent.State.Load(),
		}, nil
	})
}
