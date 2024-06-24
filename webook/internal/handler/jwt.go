package handler

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type jwtHandler struct {
	accessKey  []byte
	refreshKey []byte
}

func newJwtHandler() jwtHandler {
	return jwtHandler{
		accessKey:  []byte("fsAck3=%n*&*6XxbCd5ksXGjLHZT2fXc"),
		refreshKey: []byte("ANG0SxUAgrAsDcji1pAUYJglAVuBS0Qa"),
	}
}

type UserClaims struct {
	jwt.RegisteredClaims
	UID       int64  `json:"uid"`
	UserAgent string `json:"userAgent"`
}

type RefreshClaims struct {
	jwt.RegisteredClaims
	UID int64 `json:"uid"`
}

func (j *jwtHandler) setJWTToken(c *gin.Context, uid int64) error {
	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().AddDate(0, 0, 7)),
		},
		UID: uid,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(j.accessKey)
	if err != nil {
		return err
	}

	c.Header("x-jwt-token", tokenStr)
	return nil
}

func (j *jwtHandler) refreshToken(c *gin.Context, uid int64) error {

	claims := RefreshClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().AddDate(0, 1, 0)),
		},
		UID: uid,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(j.accessKey)
	if err != nil {
		return err
	}

	c.Header("x-refresh-token", tokenStr)
	return nil
}

func ExtractToken(c *gin.Context) string {
	header := c.GetHeader("Authorization")
	segment := strings.Split(header, " ")
	if len(segment) != 2 {
		return ""
	}

	return segment[0]
}
