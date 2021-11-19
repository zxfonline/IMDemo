package config

import (
	"github.com/sirupsen/logrus"
	baselog "github.com/zxfonline/IMDemo/core/log"
	caller "github.com/zxfonline/IMDemo/core/logrus-hook-caller"
)

func newFileCallerHook() logrus.Hook {
	if IsDebug() {
		hook_caller := caller.NewHook(&caller.CallerHookOptions{
			//DisabledField: true,
			//EnableFile:    true,
			//EnableLine:    true,
			Flags: baselog.DumpFlags,
		})
		return hook_caller
	} else {
		hook_caller := caller.NewHook(&caller.CallerHookOptions{})
		return hook_caller
	}
}
