package local

import (
	"context"
	"errors"
	"fmt"

	"github.com/lcsin/goprojets/webook/pkg/ratelimiter"
)

type Service struct {
	// 嵌入限流接口
	limiter ratelimiter.Limiter
}

func NewService(limiter ratelimiter.Limiter) *Service {
	return &Service{
		limiter: limiter,
	}
}

func (s *Service) Send(ctx context.Context, tplId string, args []string, number ...string) error {
	// 在调用第三方服务发送短信之前进行限流
	limit, err := s.limiter.Limit(ctx, "sms:local")
	if err != nil {
		return errors.New("sms limiter error")
	}
	if limit {
		return errors.New("too many request")
	}
	// 打印验证码模拟调用第三方服务发送短信
	fmt.Println(args)
	return nil
}
