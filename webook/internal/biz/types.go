package biz

import "github.com/golang-jwt/jwt/v5"

const (
	JwtKey = "fsAck3=%n*&*6XxbCd5ksXGjLHZT2fXc"
)

type UserClaims struct {
	jwt.RegisteredClaims
	UID int64 `json:"uid"`
}
