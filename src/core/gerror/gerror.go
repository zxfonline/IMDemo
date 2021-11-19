package gerror

import (
	"errors"
	"fmt"
	"strconv"
)

type ErrorType int32

func (x ErrorType) String() string {
	return EnumName(ErrorType_name, int32(x))
}

func (x ErrorType) Value() int32 {
	return int32(x)
}

// EnumName is a helper function to simplify printing enums by name.  Given an enum map and a value, it returns a useful string.
func EnumName(m map[int32]string, v int32) string {
	s, ok := m[v]
	if ok {
		return s
	}
	return strconv.Itoa(int(v))
}

var ErrorType_name = map[int32]string{
	0:      "OK",
	100000: "CUSTOM_ERROR",
	100001: "SERVER_INTERNAL_ERROR",
	100002: "SERVER_CDATA_ERROR",
	100003: "SERVER_CMSG_ERROR",
	100004: "SERVER_FILE_NOT_FOUND",
	100005: "SERVER_ACCESS_REFUSED",
	200000: "CLIENT_TIMEOUT",
	200001: "CLIENT_IO_ERROR",
}

const (
	//操作成功的返回码常量
	OK ErrorType = 0

	//通用错误
	CUSTOM_ERROR ErrorType = 100000
	//服务器端一般错误
	SERVER_INTERNAL_ERROR ErrorType = 100001
	//读取客户端发送的数据异常
	SERVER_CDATA_ERROR ErrorType = 100002
	//处理客户端发送的数据异常
	SERVER_CMSG_ERROR ErrorType = 100003
	//服务器端的文件没有找到
	SERVER_FILE_NOT_FOUND ErrorType = 100004
	//服务器端的访问被拒绝
	SERVER_ACCESS_REFUSED ErrorType = 100005

	//客户端访问超时
	CLIENT_TIMEOUT ErrorType = 200000
	//客户端IO错误
	CLIENT_IO_ERROR ErrorType = 200001
)

//判断错误类型是否是允许的错误范围 true=消息内部定义的错误，false=系统级别错误，需要关闭连接
func IsCustomError(e ErrorType) bool {
	return e < CUSTOM_ERROR
}

//系统常规错误
type SysError struct {
	Code    ErrorType   `json:"code"`
	Content string      `json:"message,omitempty"`
	Cause   error       `json:"-"`
	Data    interface{} `json:"data,omitempty"`
}

func (err *SysError) Error() string {
	if err.Cause == nil {
		return fmt.Sprintf("code:%d,message:%s", err.Code, err.Content)
	} else if err.Content != "" {
		return fmt.Sprintf("code:%d,message:%s,cause:%s", err.Code, err.Content, err.Cause.Error())
	} else {
		return fmt.Sprintf("code:%d,cause:%s", err.Code, err.Cause.Error())
	}
}

func New(code ErrorType, err error) *SysError {
	if err == nil {
		return &SysError{Code: code}
	}
	return &SysError{Code: code, Content: err.Error()}
}
func NewError(code ErrorType, err string) *SysError {
	return &SysError{Code: code, Content: err}
}
func NewErrorByCode(code ErrorType) *SysError {
	return &SysError{Code: code, Content: code.String()}
}

var RecoverPanicToErr = true

func PanicValToErr(panicVal interface{}, err *error) {
	if panicVal == nil {
		return
	}
	// case nil
	switch xerr := panicVal.(type) {
	case error:
		*err = xerr
	case string:
		*err = errors.New(xerr)
	default:
		*err = fmt.Errorf("%v", panicVal)
	}
}

func PanicToErr(err *error) {
	if RecoverPanicToErr {
		if x := recover(); x != nil {
			//debug.PrintStack()
			PanicValToErr(x, err)
		}
	}
}
