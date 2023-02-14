package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"image/color"
	"liujun/Time_Micro_GateWay/server/models"
	pb "liujun/Time_Micro_GateWay/server/proto"
	"liujun/Time_Micro_GateWay/server/utils"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/afocus/captcha"
	"github.com/astaxie/beego/validation"
)

type UserHandler struct {
}

func (uh *UserHandler) UserRegister(ctx context.Context, in *pb.RegisterRequest, out *pb.RegisterResponse) error {
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
			fmt.Println(err)
			out.Code = 1
			out.Info = fmt.Sprintf("%v", err)
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

func (uh *UserHandler) UserLogin(ctx context.Context, in *pb.LoginRequest, out *pb.LoiginResponse) error {
	//从redis获取图片验证码
	redis_code, err := models.RED.Get(context.Background(), "img_code").Result()
	if err == redis.Nil {
		out.Code = 1
		out.Info = "验证码已过期"
		return nil
	}
	if redis_code != in.Captcha {
		out.Code = 1
		out.Info = "验证码错误"
		return nil
	}
	user := models.User{}
	models.DB.Table("user").Select("user.id,user.username,user.password").Where("username = ?", in.Username).First(&user)
	if user.Username == "" {
		out.Code = 1
		out.Info = "用户名不存在"
		return nil
	}
	if utils.GetMd5(in.Password) != user.Password {
		out.Code = 1
		out.Info = "用户名或者密码错误"
		return nil
	}
	token, err := utils.GenerateToken(int64(user.ID))
	if err != nil {
		out.Code = 1
		out.Info = fmt.Sprintf("%v", err)
		return nil
	}
	fmt.Println("len(token)", len(token))
	out.Code = 0
	out.Info = "登录成功"
	out.Tokent = token
	return nil
}

func (h *UserHandler) GetCaptcha(ctx context.Context, in *pb.CaptchaRequest, out *pb.CaptchaResponse) error {
	cap := captcha.New()
	// 可以设置多个字体 或使用cap.AddFont("xx.ttf")追加
	cap.SetFont("./asset/font/comictype.ttf")
	// 设置验证码大小
	cap.SetSize(128, 64)
	// 设置干扰强度
	cap.SetDisturbance(captcha.MEDIUM)
	// 设置前景色 可以多个 随机替换文字颜色 默认黑色
	cap.SetFrontColor(color.RGBA{0, 0, 0, 0})
	// 设置背景色 可以多个 随机替换背景色 默认白色
	cap.SetBkgColor(color.RGBA{255, 255, 255, 255})
	img, str := cap.Create(4, captcha.NUM)
	img_bytes, _ := json.Marshal(img)
	out.Code = str
	out.Img = img_bytes
	models.RED.SetEX(context.TODO(), "img_code", str, time.Second*60)
	return nil
}

// 实例化一个handler
func NewUserHandlerService() *UserHandler {
	handler := new(UserHandler)
	return handler
}
