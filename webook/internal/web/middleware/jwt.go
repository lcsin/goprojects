package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lcsin/gopocket/util/ginx"
	"github.com/lcsin/goprojets/webook/internal/biz"
)

func Jwt() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		// 不需要登录校验
		if path == "/api/v1/users/signup" || path == "/api/v1/users/login" || path == "/api/v1/ping" ||
			path == "/api/v1/users/login/sms/code/send" || path == "/api/v1/users/login/sms" {
			return
		}

		header := c.GetHeader("Authorization")
		if header == "" {
			ginx.ResponseError(c, ginx.ErrUnauthorized)
			c.Abort()
			return
		}

		segment := strings.Split(header, " ")
		if len(segment) != 2 {
			ginx.ResponseError(c, ginx.ErrUnauthorized)
			c.Abort()
			return
		}
		tokenStr := segment[1]
		var claims biz.UserClaims
		token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(biz.JwtKey), nil
		})
		if err != nil {
			ginx.ResponseError(c, ginx.ErrUnauthorized)
			c.Abort()
			return
		}
		if !token.Valid || claims.UserAgent != c.GetHeader("User-Agent") {
			ginx.ResponseError(c, ginx.ErrUnauthorized)
			c.Abort()
			return
		}

		c.Set("uid", claims.UID)
	}
}
