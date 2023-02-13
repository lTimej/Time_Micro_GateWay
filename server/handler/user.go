package handler

import (
	"context"
	pb "liujun/Time_Micro_GateWay/server/proto"
)

type UserHandler struct {
}

func (uh *UserHandler) UserService(ctx context.Context, in *pb.TestRequest, out *pb.TestResponse) error {

	return nil
}

func NewUserHandlerService() *UserHandler {
	handler := new(UserHandler)

	return handler
}
