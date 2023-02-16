package handler

import (
	"context"
	pb "liujun/Time_Micro_GateWay/mail/proto"
)

type MailHandler struct {
}

func (mh *MailHandler) TestMail(ctx context.Context, in *pb.MailTestRequest, out *pb.MailTestResponse) error {
	return nil
}

func NewMailHandlerService() *MailHandler {
	handler := new(MailHandler)
	return handler
}
