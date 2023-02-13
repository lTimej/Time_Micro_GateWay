package handler

import (
	"fmt"
	pb "liujun/Time_Micro_GateWay/proto"

	"context"
	"liujun/Time_Micro_GateWay/common"

	"github.com/asim/go-micro/plugins/registry/consul/v4"
	"github.com/gin-gonic/gin"
	"go-micro.dev/v4"
)

func getClient() pb.UserService {
	consulRegistry := consul.NewRegistry()
	srv := micro.NewService(
		micro.Registry(consulRegistry))
	return pb.NewUserService(common.ServiceName, srv.Client())
}

func Index(c *gin.Context) {
	client := getClient()
	resp, err := client.UserService(context.TODO(), &pb.UserRequest{Id: 10})
	if err != nil {
		fmt.Println(err)
		c.JSON(200, "错误")
		return
	}
	c.JSON(200, resp.Msg)
}
