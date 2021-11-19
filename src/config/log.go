package config

import (
	"io"
	"log"
	"os"
	"regexp"

	caller "github.com/zxfonline/IMDemo/core/logrus-hook-caller"

	"github.com/sirupsen/logrus"
	baselog "github.com/zxfonline/IMDemo/core/log"
)

var MiscRegexp = regexp.MustCompile(`core(@.*)?/log/.*.go`)

func initLogger() {
	baselog.SuffixesToIgnoreArray = append(baselog.SuffixesToIgnoreArray, MiscRegexp)
	caller.SuffixesToIgnoreArray = append(caller.SuffixesToIgnoreArray, MiscRegexp)
	// Log as JSON instead of the default ASCII formatter.
	if Conf.LogFmt == "text" {
		//logrus.SetFormatter(&logrus.TextFormatter{DisableColors: true})
		logrus.SetFormatter(&logrus.TextFormatter{})
	} else {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	if IsDebug() {
		baselog.DumpFlags = log.Llongfile
		baselog.SetOutput(io.MultiWriter(newLogWriter(), os.Stdout))
	} else {
		baselog.DumpFlags = log.Lshortfile
		baselog.SetOutput(newLogWriter())
	}

	// Only log the warning severity or above.
	if level, ok := logLevels[Conf.LogLevel]; ok {
		baselog.Logger.SetLevel(level)
	} else {
		baselog.Logger.SetLevel(logrus.InfoLevel)
	}
	//add log hook
	baselog.Logger.AddHook(newFileCallerHook())
}
