package setup

import (
	"context"
	"fmt"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/zxfonline/IMDemo/clientctl"
	"github.com/zxfonline/IMDemo/config"
	"github.com/zxfonline/IMDemo/core/fileutil"
	"github.com/zxfonline/IMDemo/core/log"
	"github.com/zxfonline/IMDemo/core/web"
	"github.com/zxfonline/IMDemo/service"
)

func Setup() {
	initEnv()
	startService()
}
func initEnv() {
	fileutil.SetOSEnv("GOTRACEBACK", "crash")
}

func startService() {
	wg := &sync.WaitGroup{}
	ctx, quitF := context.WithCancel(context.Background())
	log.Info("starting")
	defer exitHandler(ctx, quitF, wg)
	logrus.RegisterExitHandler(func() {
		exitHandler(ctx, quitF, wg)
	})
	clientctl.StartServer(ctx, wg, config.Conf.RoomSize, config.Conf.ChatCashSize)
	webServer := startHttp(config.Conf.HttpAddr)
	service.RegisterHandlers(ctx, wg, webServer)
	defer webServer.Close()

	signal.Notify(config.ExitSignal, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	log.Info("started")
	for {
		msg := <-config.ExitSignal
		log.Infof("signal:%v", msg)
		switch msg {
		case syscall.SIGHUP:
		case syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT:
			return
		}
	}
}

//创建http服务器
func startHttp(address string) *web.Server {
	//开启http服务
	httpRoot, err := fileutil.FindFullPathPath(fileutil.PathJoin(config.Conf.Runtime, "www/static"))
	if err != nil {
		panic(err)
	}
	web.CopyRequestBody = true
	web.IndentJson = false
	httpCfg := web.NewServerConfig(web.SetMaxHeaderBytes(1<<19),
		web.SetStaticDir(fileutil.TransPath(httpRoot)))
	sev := web.NewServer(web.SetServerConfig(httpCfg))
	i := strings.Index(address, ":")
	if i < 0 {
		panic(fmt.Errorf("Bad HttpAddr: %s", address))
	}
	address = address[i:]

	sev.RunMux("/", address)
	return sev
}

var exitOnce sync.Once

//关服流程
func exitHandler(ctx context.Context, cancel context.CancelFunc, wg *sync.WaitGroup) {
	exitOnce.Do(func() {
		log.Info("closing")
		//关服收尾流程开始
		cancel()
		//等待数据操作相关的线程完成任务
		wg.Wait()
		log.Info("closed")
	})
}
