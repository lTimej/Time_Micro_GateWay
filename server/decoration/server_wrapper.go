package decoration

import (
	"context"
	"log"

	"go-micro.dev/v4/server"
)

// 服务端装饰器
func ServerWrapper() server.HandlerWrapper {
	return func(h server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			log.Println("Server wrapper!!!")
			return h(ctx, req, rsp)
		}
	}
}
