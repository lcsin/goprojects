package ioc

import (
	"time"

	"github.com/lcsin/goprojets/webook/internal/service/sms"
	"github.com/lcsin/goprojets/webook/internal/service/sms/local"
	"github.com/lcsin/goprojets/webook/pkg/ratelimiter"
	"github.com/redis/go-redis/v9"
)

// InitSMSService 初始化短信服务
func InitSMSService(cmd redis.Cmdable) sms.Service {
	return local.NewService(ratelimiter.NewRedisSlideWindowLimiter(cmd, time.Second, 100))
}
