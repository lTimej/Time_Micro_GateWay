package utils

import (
	"fmt"
	pb "liujun/Time_Micro_GateWay/proto"

	"github.com/afex/hystrix-go/hystrix"
)

func HystrixLimit(request_name, endpoint string, action func() error) error {
	fmt.Println("=================", request_name, "=================")
	var config hystrix.CommandConfig
	switch request_name {
	case "UserService.UserLogin":
		config = hystrix.CommandConfig{
			Timeout:                3000,
			MaxConcurrentRequests:  2,
			ErrorPercentThreshold:  25,
			RequestVolumeThreshold: 2,
			SleepWindow:            3000,
		}
	case "UserService.UserRegister":
		config = hystrix.CommandConfig{
			Timeout:                5000,
			MaxConcurrentRequests:  10,
			ErrorPercentThreshold:  25,
			RequestVolumeThreshold: 10,
			SleepWindow:            2000,
		}
	case "UserService.GetCaptcha":
		config = hystrix.CommandConfig{
			Timeout:                1000,
			MaxConcurrentRequests:  2,
			ErrorPercentThreshold:  25,
			RequestVolumeThreshold: 2,
			SleepWindow:            2000,
		}
	default:
		config = hystrix.CommandConfig{
			Timeout:                5000,
			MaxConcurrentRequests:  10000,
			ErrorPercentThreshold:  25,
			RequestVolumeThreshold: 100,
			SleepWindow:            5000,
		}
	}
	hystrix.ConfigureCommand(request_name, config)
	return hystrix.Do(request_name, func() error {
		return action()
	}, func(err error) error {
		return err
	})
}

func HystrixFallback(rsp interface{}) {
	switch t := rsp.(type) {

	case *pb.RegisterResponse:
		t.Code = 2
		t.Info = "please wait to try"
	case *pb.LoiginResponse:
		t.Code = 2
		t.Info = "please wait to try"
	case *pb.TestResponse:
		t.Msg = "is limit"
	default:
	}
}
