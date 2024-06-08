package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lcsin/goprojets/webook/internal/service"
	"github.com/lcsin/goprojets/webook/internal/web/middleware"
)

func RegisterRoutes(us *service.UserService) *gin.Engine {
	r := gin.Default()

	// 将用户的登录信息存储在cookie
	//store := cookie.NewStore([]byte("secret"))
	// 设置中间件
	r.Use(middleware.CORS(), middleware.Jwt())

	v1 := r.Group("/api/v1")
	v1.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// user routes
	NewUserHandler(us).RegisterRoutes(v1)

	return r
}
