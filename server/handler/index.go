package handler

import (
	"context"
	"fmt"

	pb "liujun/Time_Micro_GateWay/server/proto"
)

func (uh *UserHandler) TestUser(ctx context.Context, in *pb.TestRequest, out *pb.TestResponse) error {
	//给out赋值
	out.Msg = fmt.Sprintf("%s===%d", "hahaha", in.Id)
	return nil
}
