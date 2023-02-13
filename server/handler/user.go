package handler

import (
	"context"
	"fmt"
	"liujun/Time_Micro_GateWay/server/models"
	pb "liujun/Time_Micro_GateWay/server/proto"
	"liujun/Time_Micro_GateWay/server/utils"

	"github.com/astaxie/beego/validation"
)

type UserHandler struct {
}

func (uh *UserHandler) UserRegister(ctx context.Context, in *pb.RegisterRequest, out *pb.RegisterResponse) error {
	models.DB.AutoMigrate(&models.User{})
	user := models.User{
		Username:   in.Username,
		Password:   in.Password,
		Email:      in.Email,
		Phone:      in.Phone,
		RePassword: in.RePassword,
	}
	valid := validation.Validation{}
	b, err := valid.Valid(&user)
	if err != nil {
		return err
	}
	if !b {
		for _, err = range valid.Errors {
			out.Code = 1
			out.Info = "数据为空"
			return nil
		}
	}
	if in.Password != in.RePassword {
		out.Code = 1
		out.Info = "密码不一致"
		return nil
	}
	var count int64
	models.DB.Table("user").Where("username = ?", in.Username).Count(&count)
	if count > 0 {
		out.Code = 1
		out.Info = "用户名已存在"
		return nil
	}
	models.DB.Table("user").Where("email = ?", in.Email).Count(&count)
	if count > 0 {
		out.Code = 1
		out.Info = "邮箱已注册"
		return nil
	}
	models.DB.Table("user").Where("phone = ?", in.Phone).Count(&count)
	if count > 0 {
		out.Code = 1
		out.Info = "手机号已注册"
		return nil
	}
	user.Password = utils.GetMd5(user.Password)
	fmt.Println(user.Password)
	if err := models.DB.Create(&user).Error; err != nil {
		return err
	}
	out.Code = 0
	out.Info = ""
	return nil
}

// 实例化一个handler
func NewUserHandlerService() *UserHandler {
	handler := new(UserHandler)
	return handler
}
