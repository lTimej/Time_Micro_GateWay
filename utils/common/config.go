package common

import (
	"github.com/astaxie/beego/config"
	"github.com/sirupsen/logrus"
)

var (
	Config config.Configer
	err    error
)

const (
	DefaultUserRoleId = 1
)

func init() {
	Config, err = config.NewConfig("ini", "common/conf.config")
	if err != nil {
		logrus.Println("配置错误")
	}
}
