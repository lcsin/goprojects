package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lcsin/goprojets/webook/internal/service"
	"github.com/lcsin/goprojets/webook/internal/web/middleware"
)

func RegisterRoutes(us *service.UserService, cs *service.CodeService) *gin.Engine {
	r := gin.Default()

	// 将用户的登录信息存储在cookie
	//store := cookie.NewStore([]byte("secret"))
	// 设置中间件
	r.Use(
		middleware.CORS(),
		middleware.NewJwtBuilder().IgnorePaths(
			"/api/v1/ping",
			"/api/v1/users/signup",
			"/api/v1/users/login",
			"/api/v1/users/login/sms/code/send",
			"/api/v1/users/login/sms",
		).Build(),
	)

	v1 := r.Group("/api/v1")
	v1.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// user routes
	NewUserHandler(us, cs).RegisterRoutes(v1)

	return r
}
