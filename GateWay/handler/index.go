package handler

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	pb "liujun/Time_Micro_GateWay/proto"
)

type Test struct{}

func (e *Test) UserService(ctx context.Context, in *pb.UserRequest, out *pb.UserResponse) error {
	fmt.Println("hahah")
	return nil
}

func Index(c *gin.Context) {
	c.JSON(200, "hello")
}
