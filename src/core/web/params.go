package web

import (
	"fmt"
)

//没找到并且没有默认值则 panic抛错
func (ctx *Context) ParamStr(key string, def ...string) string {
	if v, ok := ctx.Params[key]; ok {
		return v
	} else if len(def) > 0 {
		return def[0]
	}
	panic(fmt.Errorf("no found param:%v", key))
}

// Param returns router param by a given key.没找到并且没有默认值则返回""字符串
func (ctx *Context) Param(key string, def ...string) string {
	if v, ok := ctx.Params[key]; ok {
		return v
	}
	var defv string
	if len(def) > 0 {
		defv = def[0]
	}
	return defv
}
