package router

import (
	"liujun/Time_Micro_GateWay/handler"

	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	router := gin.Default()
	//测试例子
	router.GET("/index", handler.Index)
	//注册
	router.POST("/register", handler.UserRegister)
	//获取验证吗
	router.GET("/getcaptcha", handler.GetCaptcha)
	return router
}
