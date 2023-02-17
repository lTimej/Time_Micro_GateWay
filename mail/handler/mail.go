package handler

import (
	"context"
	pb "liujun/Time_Micro_GateWay/mail/proto"
	"log"
)

type MailHandler struct {
}

func (mh *MailHandler) TestMail(ctx context.Context, in *pb.MailTestRequest, out *pb.MailTestResponse) error {
	log.Println("test mail send... ...")
	out.Msg = "success"
	return nil
}

func NewMailHandlerService() *MailHandler {
	handler := new(MailHandler)
	return handler
}
