package decoration

import (
	"context"
	"errors"
	"fmt"
	"liujun/Time_Micro_GateWay/server/utils"
	"log"

	"go-micro.dev/v4/server"
	// "google.golang.org/grpc/metadata"
)

// 服务端装饰器
func ServerWrapper() server.HandlerWrapper {
	return func(h server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			log.Println(req.Endpoint())
			if req.Endpoint() == "UserService.GetCaptcha" || req.Endpoint() == "UserService.UserLogin" {
				return h(ctx, req, rsp)
			}
			header := req.Header()
			if header == nil {
				log.Println("认证失败")
				return errors.New("认证失败")
			}
			// var token string
			// if md, ok := metadata.FromIncomingContext(ctx); ok {
			// 	for key, val := range md {
			// 		if key == "Authorization" {
			// 			token = val[0]
			// 		}
			// 	}
			// } else {
			// 	log.Println("认证失败", md)
			// 	return errors.New("认证失败")
			// }
			fmt.Println("------------>", header)
			token := header["Authorization"]
			if token == "" {
				log.Println("认证失败，缺少token")
				return errors.New("认证失败，缺少token")
			}
			fmt.Printf("%T\n", token)
			err := utils.VerifyToken(token)
			if err != nil {
				return err
			}
			return h(ctx, req, rsp)
		}
	}
}
