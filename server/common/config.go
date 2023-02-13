package common

import (
	"github.com/astaxie/beego/config"
	"github.com/sirupsen/logrus"
)

var (
	Config config.Configer
	err    error
)

func init() {
	Config, err = config.NewConfig("ini", "./conf/conf.config")
	if err != nil {
		logrus.Println("配置错误")
	}
}
