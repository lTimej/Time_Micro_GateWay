package handler

import (
	"context"
	"encoding/json"
	"fmt"
	pb "liujun/Time_Micro_GateWay/proto"
	"liujun/Time_Micro_GateWay/utils"
	"log"

	"github.com/afocus/captcha"
	"github.com/gin-gonic/gin"
)

func UserRegister(c *gin.Context) {
	data, _ := c.GetRawData()
	var args map[string]interface{}
	json.Unmarshal(data, &args)
	fmt.Println(args)
	client := getClient()
	resp, err := client.UserRegister(context.Background(), &pb.RegisterRequest{
		Username:   utils.GetString(args["username"]),
		Password:   utils.GetString(args["password"]),
		Email:      utils.GetString(args["email"]),
		Phone:      utils.GetString(args["phone"]),
		RePassword: utils.GetString(args["re_password"]),
	})
	if err != nil {
		log.Println("获取数据错误err:", err)
		c.JSON(200, gin.H{"code": 1, "info": "获取数据错误"})
		return
	}
	res := make(map[string]interface{})
	res["code"] = resp.Code
	res["info"] = resp.Info
	c.JSON(200, res)
}

func GetCaptcha(c *gin.Context) {
	client := getClient()
	resp, err := client.GetCaptcha(context.Background(), &pb.CaptchaRequest{})
	if err != nil {
		log.Println("获取图片验证码错误,err:", err)
		c.JSON(200, gin.H{"code": 1, "info": "获取图片验证码错误"})
		return
	}
	var img captcha.Image
	json.Unmarshal(resp.Img, &img)
	c.JSON(200, img)
}
