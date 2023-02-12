package common

import (
	"github.com/astaxie/beego/config"
	"github.com/sirupsen/logrus"
)

var (
	ConsulConfig config.Configer
	err          error
)

func init() {
	ConsulConfig, err = config.NewConfig("ini", "./conf/gateway.config")
	if err != nil {
		logrus.Println("配置错误")
	}
}
