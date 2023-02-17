package decoration

import (
	"context"
	"liujun/Time_Micro_GateWay/utils"

	client "go-micro.dev/v4/client"
)

type CliWrapper struct {
	client.Client
}

func (c *CliWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	request_name := req.Endpoint()
	action := func() error {
		return c.Client.Call(ctx, req, rsp, opts...)
	}
	if err := utils.HystrixLimit(request_name, req.Endpoint(), action); err == nil {
		return nil
	} else {
		utils.HystrixFallback(rsp)
		return nil
	}

}

func NewClientWrapper(c client.Client) client.Client {
	return &CliWrapper{c}
}
