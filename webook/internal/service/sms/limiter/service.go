package limiter

import (
	"context"
	"errors"
	"fmt"

	"github.com/lcsin/goprojets/webook/internal/service/sms"
	"github.com/lcsin/goprojets/webook/pkg/ratelimiter"
)

var ErrLimited = errors.New("sms limited error")

// Service 装饰器模式实现短信发送服务限流
type Service struct {
	// svc 装饰sms短信发送接口
	svc sms.Service
	// limiter 装饰内容为我们的限流接口
	limiter ratelimiter.Limiter
	// key 限流对象
	key string
}

func NewService(svc sms.Service, limiter ratelimiter.Limiter, key string) *Service {
	return &Service{
		svc:     svc,
		limiter: limiter,
		key:     key,
	}
}

func (s *Service) Send(ctx context.Context, tplId string, args []string, number ...string) error {
	// 添加限流机制
	limit, err := s.limiter.Limit(ctx, s.key)
	if err != nil {
		return fmt.Errorf("sms limiter error: %v", err)
	}
	if limit {
		return ErrLimited
	}

	err = s.svc.Send(ctx, tplId, args, number...)
	// 可以在这里添加额外的机制
	return err
}
