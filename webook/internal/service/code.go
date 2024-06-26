package service

import (
	"context"
	"fmt"
	"log"
	"math/rand"

	"github.com/lcsin/goprojets/webook/internal/repository"
	"github.com/lcsin/goprojets/webook/internal/service/sms"
)

var (
	ErrCodeVerifyTooManyTimes = repository.ErrCodeVerifyTooManyTimes
	ErrCodeSendTooMnay        = repository.ErrCodeSendTooMany
)

type CodeService struct {
	repo *repository.CodeRepository
	sms  sms.Service
}

func NewCodeService(repo *repository.CodeRepository, sms sms.Service) *CodeService {
<<<<<<< HEAD
	return &CodeService{
		repo: repo,
		sms:  sms,
	}
=======
	return &CodeService{repo: repo, sms: sms}
>>>>>>> ae6108dc1c76bf8cfdb09b580d941e47dd2f2ab6
}

func (cs *CodeService) Send(ctx context.Context, biz string, phone string) error {
	// 生成验证码和保存验证码
	code := cs.generateCode()
	if err := cs.repo.Set(ctx, biz, phone, code); err != nil {
		return err
	}
	// 发送短信验证码
	return cs.sms.Send(ctx, "tplId", []string{code}, phone)
}

func (cs *CodeService) Verify(ctx context.Context, biz string, phone, inputCode string) (bool, error) {
	ok, err := cs.repo.Verify(ctx, biz, phone, inputCode)
	// 对外屏蔽了验证次数过多的错误
	if err == repository.ErrCodeVerifyTooManyTimes {
		log.Println(err)
		return false, nil
	}
	return ok, err
}

func (cs *CodeService) generateCode() string {
	code := rand.Intn(1000000)
	return fmt.Sprintf("%06d", code)
}
