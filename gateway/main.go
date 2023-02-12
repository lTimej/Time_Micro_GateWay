package main

import (
	"github.com/asim/go-micro/plugins/registry/consul/v4"
	"go-micro.dev/v4"
	"liujun/Time_Micro_GateWay/handler"
	pb "liujun/Time_Micro_GateWay/proto"
	"liujun/Time_Micro_GateWay/router"
)

func main() {
	consulRegistry := consul.NewRegistry()
	srv := micro.NewService(
		micro.Registry(consulRegistry))
	pb.NewUserService("test", srv.Client())

	r := router.Router()
	r.GET("/", handler.Index)
}
