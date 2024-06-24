package failover

import (
	"context"
	"errors"
	"log"
	"sync/atomic"

	"github.com/lcsin/goprojets/webook/internal/service/sms"
)

type Service struct {
	svcs []sms.Service
	idx  uint64
}

func NewService(svcs []sms.Service) *Service {
	return &Service{svcs: svcs}
}

func (s *Service) Send(ctx context.Context, tplId string, args []string, number ...string) error {
	idx := atomic.AddUint64(&s.idx, 1)
	length := uint64(len(s.svcs))
	for i := idx; i < idx+length; i++ {
		svc := s.svcs[i%length]
		err := svc.Send(ctx, tplId, args, number...)
		switch err {
		case nil:
			return nil
		case context.DeadlineExceeded, context.Canceled:
			// 调用者设置了超时时间或者调用者主动取消了
			return err
		}
		// 其他情况需要打印日志做好监控
		log.Println(err)
	}

	// 所有服务商发送失败大概率是你网络问题
	return errors.New("所有的短信服务商全部发送失败")
}
