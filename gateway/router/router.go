package router

import (
	"github.com/gin-gonic/gin"
	"liujun/Time_Micro_GateWay/gateway/handler"
)

func Router() *gin.Engine {
	router := gin.Default()

	router.GET("/user/test/:user_id", handler.Index)

	return router
}
