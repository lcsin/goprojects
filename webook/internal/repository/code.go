package repository

import (
	"context"

	"github.com/lcsin/goprojets/webook/internal/repository/cache"
)

var (
	ErrCodeVerifyTooManyTimes = cache.ErrCodeVerifyTooManyTimes
	ErrCodeSendTooMany        = cache.ErrCodeSendToMany
)

type ICodeRepository interface {
	Set(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz, phone, inputCode string) (bool, error)
}

type CodeRepository struct {
	cache cache.ICodeCache
}

func NewCodeRepository(cache cache.ICodeCache) ICodeRepository {
	return &CodeRepository{cache: cache}
}

func (cr *CodeRepository) Set(ctx context.Context, biz, phone, code string) error {
	return cr.cache.Set(ctx, biz, phone, code)
}

func (cr *CodeRepository) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	return cr.cache.Verify(ctx, biz, phone, inputCode)
}
