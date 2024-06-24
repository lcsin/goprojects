package ratelimiter

import "context"

type Limiter interface {
	// Limit 限流方法
	// key:限流对象, bool:是否限流, error:限流器是否有错误
	Limit(ctx context.Context, key string) (bool, error)
}
