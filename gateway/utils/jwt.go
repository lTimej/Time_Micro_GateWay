package utils

import (
	"errors"

	"github.com/golang-jwt/jwt/v4"
)

var CustomSecret = []byte("micr0_gatevvay")

type CustomClaims struct {
	UserId int64 `json:"user_id"`
	RoleId int64 `json:"role_id"`
	jwt.RegisteredClaims
}

func GetTokenData(token_str string) (int64, int64, error) {
	if token_str == "" {
		return 0, 0, errors.New("token不能为空")
	}
	token, err := jwt.ParseWithClaims(token_str, &CustomClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return CustomSecret, nil
	})
	if err != nil {
		return 0, 0, errors.New("非法token")
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims.UserId, claims.RoleId, nil
	}
	return 0, 0, err
}
