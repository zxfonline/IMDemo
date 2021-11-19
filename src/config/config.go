package config

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	"github.com/zxfonline/IMDemo/core/fileutil"
	"gopkg.in/yaml.v2"
)

const DEBUG = "DEBUG"
const RELEASE = "RELEASE"

type Config struct {
	Mode string `yaml:"mode"`
	ID   int32  `yaml:"id"`

	RoomSize     int64 `yaml:"roomSize"`
	ChatCashSize int32 `yaml:"chatCashSize"`

	Runtime string `yaml:"runtime"`
	Static  string `yaml:"static"`
	//服务器等级
	Level        int32  `yaml:"level"`
	LogLocalPath string `yaml:"logLocalPath"`
	LogLevel     string `yaml:"logLevel"`
	LogFmt       string `yaml:"logFmt"`
	HttpAddr     string `yaml:"httpAddr"`
}

var (
	Conf         Config
	_config_file = flag.String("c", "", "config filename")

	//进程退出管道
	ExitSignal = make(chan os.Signal, 1)
)

//GetConfig 获取启动配置
func GetConfig() *Config {
	return &Conf
}

// init the Conf
func init() {
	flag.Parse()
	rand.Seed(time.Now().UnixNano())
	if len(*_config_file) == 0 {
		if fullPath, err := fileutil.FindFullFilePath("runtime/config.yml"); err == nil {
			*_config_file = fullPath
		}
	}
	if data, err := ioutil.ReadFile(*_config_file); err != nil {
		panic(fmt.Errorf("load start config err:%v[default:./runtime/config.yml or ../runtime/config.yml]", err))
	} else if err = yaml.Unmarshal(data, &Conf); err != nil {
		panic(fmt.Errorf("unmarshal start config err:%v", err))
	}
	if err := checkConfig(); err != nil {
		panic(fmt.Errorf("check start config err:%v", err))
	}
	if httpUploads, err := fileutil.FindFullPathPath(Conf.Static); err != nil {
		panic(err)
	} else {
		Conf.Static = httpUploads
	}
	if err := initUUID(); err != nil {
		panic(err)
	}
	initLogger()

}

func checkConfig() error {
	if ver := Conf.Mode; ver != DEBUG && ver != RELEASE {
		return fmt.Errorf("mode must be '%s' or '%s'", DEBUG, RELEASE)
	}
	return nil
}

func IsDebug() bool {
	return Conf.Mode == DEBUG
}
