package decoration

import (
	"context"
	"errors"
	opentracing "github.com/asim/go-micro/plugins/wrapper/trace/opentracing/v4"
	"go-micro.dev/v4/server"
	"liujun/Time_Micro_GateWay/server/common"
	"liujun/Time_Micro_GateWay/server/utils"
	"log"
)

// 服务端装饰器
func ServerWrapper() server.HandlerWrapper {
	return func(h server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			if req.Endpoint() == "UserService.GetCaptcha" || req.Endpoint() == "UserService.UserLogin" || req.Endpoint() == "UserService.UserRegister" {
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

// 分布式链路追踪
func ServerTrace(service *common.MicroService) server.HandlerWrapper {
	tracer := utils.GetTracer(common.ServiceName, service.Addr, common.Config)
	tracerHandler := opentracing.NewHandlerWrapper(tracer)
	return tracerHandler
}
