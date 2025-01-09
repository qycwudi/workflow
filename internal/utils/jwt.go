package utils

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

// @secretKey: JWT 加解密密钥
// @iat: 时间戳
// @seconds: 过期时间，单位秒
// @userId: 用户id
func GenerateJwtToken(secretKey string, iat, seconds int64, userId int64) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims["userId"] = userId
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}

func GetUserId(ctx context.Context) (int64, error) {
	userId := ctx.Value("userId")
	if userId == nil {
		return -1, errors.New("user not login")
	}

	// 将json.Number类型转换为int64
	userIdNum, ok := userId.(json.Number)
	if !ok {
		return -1, errors.New("invalid user id")
	}

	userIdInt64, err := userIdNum.Int64()
	if err != nil {
		return -1, errors.New("invalid user id")
	}
	return userIdInt64, nil
}
