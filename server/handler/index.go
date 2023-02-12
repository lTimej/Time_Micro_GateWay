package handler

import (
	"context"
	"fmt"

	pb "liujun/Time_Micro_GateWay/server/proto"
)

type Test struct{}

func (e *Test) UserService(ctx context.Context, in *pb.UserRequest, out *pb.UserResponse) error {
	fmt.Println("hahah")
	return nil
}
