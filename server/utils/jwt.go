package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const TokenExpireDuration = time.Hour * 2

// CustomSecret 用于加盐的字符串
var CustomSecret = []byte("helloworld")

type CustomClaims struct {
	// 可根据需要自行添加字段
	UserId             int64 `json:"user_id"`
	jwt.StandardClaims       // 内嵌标准的声明
}

// func GenerateToken(user_id int64) (string, error) {
// 	claims := CustomClaims{
// 		user_id, // 自定义字段
// 		jwt.StandardClaims{
// 			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(),
// 			Issuer:    "microgateway", // 签发人
// 		},
// 	}
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	return token.SignedString(CustomSecret)
// }

// func VerifyToken(token_str string) error {
// 	if token_str == "" {
// 		return errors.New("token不能为空")
// 	}
// 	fmt.Println("token_str", token_str)

// 	// tokens, err := jwt.Parse(token_str, func(token *jwt.Token) (interface{}, error) {
// 	// 	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 	// 		return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
// 	// 	}
// 	// 	return []byte(CustomSecret), nil
// 	// })

// 	// if claims, ok := tokens.Claims.(jwt.MapClaims); ok && tokens.Valid {
// 	// 	if int(claims["UserId"].(float64)) <= 0 {
// 	// 		return errors.New("auth token: data is 0")
// 	// 	}
// 	// 	return nil
// 	// } else {
// 	// 	fmt.Println("auth token:", err)
// 	// 	return err
// 	// }

// 	token, err := jwt.ParseWithClaims(token_str, &CustomClaims{}, func(token *jwt.Token) (i interface{}, err error) {
// 		return CustomSecret, nil
// 	})
// 	if err != nil {
// 		log.Println("5555555555555", err)
// 		return err
// 	}
// 	if _, ok := token.Claims.(*CustomClaims); ok && token.Valid {
// 		return nil
// 	}
// 	return errors.New("非法token")
// }

func GenerateToken(user_id int64) (string, error) {
	// 创建 Claims
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * 30)), // 过期时间
		Issuer:    "Time",                                               // 签发人
	}
	// 生成token对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 生成签名字符串
	return token.SignedString(CustomSecret)
}
func VerifyToken(token_str string) error {
	fmt.Println("token_str", token_str)
	token, err := jwt.Parse(token_str, func(token *jwt.Token) (interface{}, error) {
		fmt.Println("99999988888999999999")
		return CustomSecret, nil
	})
	fmt.Println("解析后的token:", token)
	if err != nil {
		fmt.Println("解析错误,err:", err)
		return err
	}
	fmt.Println(token)
	if _, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return nil
	}
	return errors.New("非法token")
}
