package main

import (
	"fmt"
	"github.com/asim/go-micro/plugins/registry/consul/v4"
	"go-micro.dev/v4"
	"go-micro.dev/v4/logger"
	"go-micro.dev/v4/registry"
	"liujun/Time_Micro_GateWay/mail/common"
	"liujun/Time_Micro_GateWay/mail/handler"
	pb "liujun/Time_Micro_GateWay/mail/proto"
	"time"
)

func main() {
	config := common.Config
	consul_addr := config.String("consul_addr")
	consul_port, _ := config.Int("consul_addr")
	consul_registry := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			fmt.Sprintf("%s:%d", consul_addr, consul_port),
		}
	})
	srv := micro.NewService(
		micro.Name(common.ServiceName),
		micro.Version(common.Version),
		micro.Registry(consul_registry),
		micro.Address(":8081"), //Transport [http] Listening on [::]:8081
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*30))
	if err := pb.RegisterMailServiceHandler(srv.Server(), handler.NewMailHandlerService()); err != nil {
		logger.Fatal(err)
	}
	// 启动服务
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
