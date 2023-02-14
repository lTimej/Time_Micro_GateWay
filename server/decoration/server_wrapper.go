package decoration

import (
	"context"
	"errors"
	"liujun/Time_Micro_GateWay/server/utils"
	"log"

	"go-micro.dev/v4/server"
	// "google.golang.org/grpc/metadata"
)

// 服务端装饰器
func ServerWrapper() server.HandlerWrapper {
	return func(h server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			if req.Endpoint() == "UserService.GetCaptcha" || req.Endpoint() == "UserService.UserLogin" {
				return h(ctx, req, rsp)
			}
			header := req.Header()
			if header == nil {
				log.Println("认证失败")
				return errors.New("认证失败")
			}
			token := header["Authorization"]
			if token == "" {
				log.Println("认证失败，缺少token")
				return errors.New("认证失败，缺少token")
			}
			err := utils.VerifyToken(token)
			if err != nil {
				log.Println(err)
				return err
			}
			return h(ctx, req, rsp)
		}
	}
}
