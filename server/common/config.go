package common

import (
	"github.com/astaxie/beego/config"
	"github.com/sirupsen/logrus"
	"go-micro.dev/v4"
	"liujun/Time_Micro_GateWay/server/utils"
	"strings"
)

var (
	Config config.Configer
	err    error
)

const (
	DefaultUserRoleId = 1
)

// 服务
type MicroService struct {
	Service micro.Service
	Host    string
	Port    uint
	Addr    string
}

func init() {
	Config, err = config.NewConfig("ini", "./conf/conf.config")
	if err != nil {
		logrus.Println("配置错误")
	}
}

// 创建服务
func NewService() *MicroService {
	service := new(MicroService)
	return service
}

// 初始化服务器地址
func InitAddr(service *MicroService) {
	service.Addr = service.Service.Server().Options().Address
	args := strings.Split(service.Service.Server().Options().Address, ":")
	if len(args) == 2 {
		service.Host = args[0]
		service.Port = uint(utils.GetInt(args[1]))
	}
}
