package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const TokenExpire = time.Hour * 2

var CustomSecret = []byte("micr0_gatevvay")

type CustomClaims struct {
	UserId int64 `json:"user_id"`
	RoleId int64 `json:"role_id"`
	jwt.RegisteredClaims
}

func GenerateToken(user_id int64, role_id int64) (string, error) {
	claims := CustomClaims{
		user_id,
		role_id,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpire)),
			Issuer:    "micro-gateway",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(CustomSecret)
}

func VerifyToken(token_str string) error {
	if token_str == "" {
		return errors.New("token不能为空")
	}
	token, err := jwt.ParseWithClaims(token_str, &CustomClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return CustomSecret, nil
	})
	if err != nil {
		return errors.New("非法token")
	}
	if _, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return nil
	}
	return errors.New("非法token")
}
