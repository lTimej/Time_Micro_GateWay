package decoration

import (
	"context"
	"log"

	"go-micro.dev/v4/client"
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
