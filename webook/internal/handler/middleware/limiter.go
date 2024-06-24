package middleware

import (
	_ "embed"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lcsin/goprojets/webook/pkg/ratelimiter"
)

type LimiterBuilder struct {
	prefix  string
	limiter ratelimiter.Limiter
}

func NewLimiterBuilder(limiter ratelimiter.Limiter) *LimiterBuilder {
	return &LimiterBuilder{
		prefix:  "ip-limiter",
		limiter: limiter,
	}
}

func (b *LimiterBuilder) Prefix(prefix string) *LimiterBuilder {
	b.prefix = prefix
	return b
}

func (b *LimiterBuilder) limit(c *gin.Context) (bool, error) {
	key := fmt.Sprintf("%s:%s", b.prefix, c.ClientIP())
	return b.limiter.Limit(c, key)
}

func (b *LimiterBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		limited, err := b.limit(ctx)
		if err != nil {
			log.Println(err)
			// 这一步很有意思，就是如果这边出错了
			// 要怎么办？
			// 保守做法：因为借助于 Redis 来做限流，那么 Redis 崩溃了，为了防止系统崩溃，直接限流
			ctx.AbortWithStatus(http.StatusInternalServerError)
			// 激进做法：虽然 Redis 崩溃了，但是这个时候还是要尽量服务正常的用户，所以不限流
			// ctx.Next()
			return
		}
		if limited {
			log.Println(err)
			ctx.AbortWithStatus(http.StatusTooManyRequests)
			return
		}
		ctx.Next()
	}
}
