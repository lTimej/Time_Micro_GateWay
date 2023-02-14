package handler

import (
	"context"
	"github.com/asim/go-micro/plugins/registry/consul/v4"
	"github.com/gin-gonic/gin"
	"go-micro.dev/v4"
	"go-micro.dev/v4/metadata"
	"liujun/Time_Micro_GateWay/common"
	pb "liujun/Time_Micro_GateWay/proto"
	"log"
	// "google.golang.org/grpc/metadata"
	// "github.com/asim/go-micro/v4/metadata"
)

func getClient(c *gin.Context) (pb.UserService, context.Context) {
	consulRegistry := consul.NewRegistry()
	srv := micro.NewService(
		micro.Registry(consulRegistry))
	token := c.Request.Header.Get("Authorization")
	ctx := metadata.Set(context.Background(), "Authorization", token)
	return pb.NewUserService(common.ServiceName, srv.Client()), ctx
}

func Index(c *gin.Context) {
	client, ctx := getClient(c)
	resp, err := client.TestUser(ctx, &pb.TestRequest{Id: 10})
	if err != nil {
		log.Println(err)
		c.JSON(200, "错误")
		return
	}
	c.JSON(200, resp.Msg)
}
