package aliyun

import "context"

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s Service) Send(ctx context.Context, tplId string, args []string, number ...string) {
	//TODO implement me
	panic("implement me")
}
