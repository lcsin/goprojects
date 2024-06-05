package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lcsin/goprojets/webook/internal/service"
	"github.com/lcsin/goprojets/webook/internal/web/middleware"
)

func RegisterRoutes(us *service.UserService) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CORS())

	v1 := r.Group("/api/v1")
	v1.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// user routes
	NewUserHandler(us).RegisterRoutes(v1)

	return r
}
