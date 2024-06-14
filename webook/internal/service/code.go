package service

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/lcsin/goprojets/webook/internal/repository"
	"github.com/lcsin/goprojets/webook/internal/service/sms"
)

type CodeService struct {
	repo *repository.CodeRepository
	sms  sms.Service
}

func (cs *CodeService) Send(ctx context.Context, biz string, phone string) error {
	// 生成验证码和保存验证码
	code := cs.GenerateCode()
	if err := cs.repo.Set(ctx, biz, phone, code); err != nil {
		return err
	}
	// 发送短信验证码
	return cs.sms.Send(ctx, "tplId", []string{code}, phone)
}

func (cs *CodeService) Verify(ctx context.Context, biz string, code, phone string) (bool, error) {

	return true, nil
}

func (cs *CodeService) GenerateCode() string {
	code := rand.Intn(1000000)
	return fmt.Sprintf("%06d", code)
}
