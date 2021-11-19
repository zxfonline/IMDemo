package config

import (
	"fmt"
	"io"
	"os"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/zxfonline/IMDemo/core/fileutil"
)

func newLogWriter() io.Writer {
	logName := fileutil.PathJoin(Conf.LogLocalPath, fmt.Sprintf("chat%d.log", Conf.ID))
	if f, err := fileutil.OpenFile(logName+"_temp", fileutil.DefaultFileFlag, fileutil.DefaultFileMode); err != nil {
		panic(fmt.Errorf("init log tmp file err:%v", err))
	} else if err := f.Close(); err != nil {
		panic(fmt.Errorf("close log tmp file err:%v", err))
	} else if err := os.Remove(f.Name()); err != nil {
		panic(fmt.Errorf("remove log tmp file err:%v", err))
	}
	//.%Y%m%d%H%M%S
	writer, err := rotatelogs.New(logName+".%Y%m%d",
		// WithRotationTime设置日志分割的时间
		rotatelogs.WithRotationTime(time.Hour*24),

		// WithMaxAge和WithRotationCount二者只能设置一个，
		// WithMaxAge设置文件清理前的最长保存时间，
		// WithRotationCount设置文件清理前最多保存的个数。
		//rotatelogs.WithMaxAge(time.Hour*24*365),
		rotatelogs.WithRotationCount(365),
	)
	if err != nil {
		panic(fmt.Errorf("config local file system for logger error: %v", err))
	}
	return writer
}
