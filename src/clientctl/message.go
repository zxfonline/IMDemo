package clientctl

import (
	"encoding/json"

	"github.com/zxfonline/IMDemo/core/gerror"
)

//消息类型枚举
type RequestType uint

const (
	LoginReq RequestType = 1001 //登录 请求
	LoginAck RequestType = 1002 //登录 响应

	RoomSwitchReq RequestType = 2001 //切换房间 请求
	RoomSwitchAck RequestType = 2002 //切换房间 响应

	RoomChatReq RequestType = 3001 //发送聊天消息 请求
	RoomChatAck RequestType = 3002 //发送聊天消息 响应
	RoomChatNtf RequestType = 4001 //聊天消息 广播
)

type Response struct {
	Type RequestType `json:"type"`
	//=0:默认成功码
	//<>0其他错误码
	Code    gerror.ErrorType `json:"code"`
	CodeMsg string           `json:"message,omitempty"`
	Data    interface{}      `json:"data,omitempty"` // 数据 json
}

func (r *Response) toJson() []byte {
	if byteData, err := json.Marshal(r); err == nil {
		return byteData
	}
	return nil
}

type Request struct {
	Type RequestType `json:"type"`
	Data interface{} `json:"data"` // 数据 json
}

const (
	ERROR_IGNORE      = -1 //忽略错误操作
	ERROR_NAME_REPEAT = -2 //姓名重复
)
