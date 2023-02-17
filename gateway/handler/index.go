package handler

import (
	"context"
	"liujun/Time_Micro_GateWay/common"
	"liujun/Time_Micro_GateWay/decoration"
	pb "liujun/Time_Micro_GateWay/proto"
	"liujun/Time_Micro_GateWay/utils"
	"liujun/Time_Micro_GateWay/vendors/rbac"
	"log"

	"github.com/asim/go-micro/plugins/registry/consul/v4"
	"github.com/gin-gonic/gin"
	"go-micro.dev/v4"
	"go-micro.dev/v4/metadata"
	// "google.golang.org/grpc/metadata"
	// "github.com/asim/go-micro/v4/metadata"
)

func getClient(c *gin.Context) (pb.UserService, context.Context) {
	consulRegistry := consul.NewRegistry()
	srv := micro.NewService(
		micro.Registry(consulRegistry),
		micro.WrapClient(decoration.NewClientWrapper))
	token := c.Request.Header.Get("Authorization")
	ctx := metadata.Set(context.Background(), "Authorization", token)
	return pb.NewUserService(common.ServiceName, srv.Client()), ctx
}

func Index(c *gin.Context) {
	client, ctx := getClient(c)
	resp, err := client.TestUser(ctx, &pb.TestRequest{Id: 10})
	//token认证
	if err != nil {
		log.Println(err)
		c.JSON(200, err)
		return
	}
	//权限认证
	token := c.Request.Header.Get("Authorization")
	user_id, role_id, err := utils.GetTokenData(token)
	if err != nil {
		log.Println(err)
		c.JSON(200, err)
		return
	}
	if !rbac.RbacFilter(int(role_id), "user.TestUser") {
		c.JSON(200, "无权限访问")
		return
	}
	log.Println(user_id, role_id)
	c.JSON(200, resp.Msg)
}
