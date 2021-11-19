package config

import (
	"github.com/bwmarrin/snowflake"
)

var (
	_uuidGenerate *snowflake.Node
)

func initUUID() (err error) {
	snowflake.NodeBits = 12
	snowflake.StepBits = 10
	_uuidGenerate, err = snowflake.NewNode(int64(Conf.ID))
	return
}

//GenerateUUID 获取消息唯一id
func GenerateUUID() snowflake.ID {
	return _uuidGenerate.Generate()
}
