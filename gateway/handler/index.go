package handler

import (
	"context"
	"fmt"
	"liujun/Time_Micro_GateWay/common"
	pb "liujun/Time_Micro_GateWay/proto"
	"log"
	"os"

	"github.com/asim/go-micro/plugins/registry/consul/v4"
	"github.com/gin-gonic/gin"
	"go-micro.dev/v4"
	"go-micro.dev/v4/metadata"
	// "google.golang.org/grpc/metadata"
	// "github.com/asim/go-micro/v4/metadata"
)

func getClient() (pb.UserService, context.Context) {
	consulRegistry := consul.NewRegistry()
	srv := micro.NewService(
		micro.Registry(consulRegistry))
	file, err := os.Open("log.log")
	if err != nil {
		log.Println(err)
	}
	token := make([]byte, 1024)
	file.Read(token)
	fmt.Println("=========token:", string(token), "====================")
	// md := metadata.New(map[string]string{
	// 	"Authorization": string(token),
	// })
	// ctx := metadata.NewOutgoingContext(context.Background(), md)
	ctx := metadata.Set(context.Background(), "Authorization", string(token))
	fmt.Println(99999999)
	return pb.NewUserService(common.ServiceName, srv.Client()), ctx
}

func Index(c *gin.Context) {
	client, ctx := getClient()
	resp, err := client.TestUser(ctx, &pb.TestRequest{Id: 10})
	if err != nil {
		fmt.Println(11111, err)
		c.JSON(200, "错误")
		return
	}
	c.JSON(200, resp.Msg)
}
