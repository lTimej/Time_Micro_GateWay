package main

import (
	"fmt"
	"liujun/Time_Micro_GateWay/server/common"
	"liujun/Time_Micro_GateWay/server/decoration"
	"liujun/Time_Micro_GateWay/server/handler"
	"liujun/Time_Micro_GateWay/server/models"

	consul "github.com/asim/go-micro/plugins/registry/consul/v4"
	"go-micro.dev/v4"
	"go-micro.dev/v4/logger"
	"go-micro.dev/v4/registry"

	pb "liujun/Time_Micro_GateWay/server/proto"
	"time"
)

func main() {
	//配置中心
	addr := common.Config.String("consul_addr")
	port, _ := common.Config.Int("consul_port")
	//注册到consul
	consul_registry := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			fmt.Sprintf("%s:%d", addr, port),
		}
	})
	models.DBMigrate()
	models.DBInit()
	srv := micro.NewService(
		micro.Name(common.ServiceName),
		micro.Version(common.Version),
		micro.Registry(consul_registry),
		micro.Address(":8081"), //Transport [http] Listening on [::]:8081
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*30),
		micro.WrapHandler(decoration.ServerWrapper()),
		micro.WrapClient(decoration.NewClientWrapper),
	)

	//注册handler,将Test实例注册到服务，供客户端回调时调用
	if err := pb.RegisterUserServiceHandler(srv.Server(), handler.NewUserHandlerService()); err != nil {
		logger.Fatal(err)
	}
	// 启动服务
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
