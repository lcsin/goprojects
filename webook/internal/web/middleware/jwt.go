package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lcsin/goprojets/webook/internal/config"
)

func Jwt() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		// 不需要登录校验
		if path == "/api/v1/users/signup" || path == "/api/v1/users/login" {
			return
		}

		header := c.GetHeader("Authorization")
		if header == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		segment := strings.Split(header, " ")
		if len(segment) != 2 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenStr := segment[1]
		var claims config.UserClaims
		token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.JwtKey), nil
		})
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if !token.Valid {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("uid", claims.UID)
	}
}
