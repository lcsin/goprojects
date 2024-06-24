package ioc

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lcsin/goprojets/webook/internal/handler"
	"github.com/lcsin/goprojets/webook/internal/handler/middleware"
	"github.com/lcsin/goprojets/webook/pkg/ratelimiter"
	"github.com/redis/go-redis/v9"
)

func InitWebServer(middlewares []gin.HandlerFunc, userHandler *handler.UserHandler) *gin.Engine {
	r := gin.Default()
	r.Use(middlewares...)
	v1 := r.Group("/api/v1")

	userHandler.RegisterRoutes(v1)

	return r
}

func InitMiddlewares(redisClient redis.Cmdable) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		// 跨域中间件
		middleware.CORS(),
		// JWT中间件
		middleware.Jwt(),
		// 限流中间件
		middleware.NewLimiterBuilder(ratelimiter.NewRedisSlideWindowLimiter(redisClient, time.Second, 100)).Build(),
	}
}
