package middleware

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		// 不需要登录校验
		if path == "/api/v1/users/signup" || path == "/api/v1/users/login" {
			return
		}
		sess := sessions.Default(c)
		if sess.Get("uid") == nil {
			// 中断，不要往后执行，也就是不要执行后面的业务逻辑
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}
