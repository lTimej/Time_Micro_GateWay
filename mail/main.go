package main

import (
	"fmt"
	"go-micro.dev/v4/registry"
	"liujun/Time_Micro_GateWay/mail/common"
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
}
