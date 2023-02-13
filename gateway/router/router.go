package router

import (
	"liujun/Time_Micro_GateWay/handler"

	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	router := gin.Default()

	router.GET("/index", handler.Index)

	return router
}
