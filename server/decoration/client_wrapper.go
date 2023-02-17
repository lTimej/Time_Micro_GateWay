package decoration

import (
	"context"
	"github.com/asim/go-micro/plugins/wrapper/trace/opentracing/v4"
	"go-micro.dev/v4/client"
	"liujun/Time_Micro_GateWay/server/common"
	"liujun/Time_Micro_GateWay/server/utils"
	"log"
)

type ClientWrapper struct {
	client.Client
}

func (cw *ClientWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	log.Println("client wrapper!!!")
	return cw.Client.Call(ctx, req, rsp, opts...)
}

func NewClientWrapper(c client.Client) client.Client {
	return &ClientWrapper{c}
}
func ClientTrace(service *common.MicroService) client.Wrapper {
	tracer := utils.GetTracer(common.ServiceName, service.Addr, common.Config)
	tracerHandler := opentracing.NewClientWrapper(tracer)
	return tracerHandler
}
