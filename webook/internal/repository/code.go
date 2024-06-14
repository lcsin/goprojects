package repository

import (
	"context"

	"github.com/lcsin/goprojets/webook/internal/repository/cache"
)

var ErrCodeVerifyTooMany = cache.ErrSendCodeToMany

type CodeRepository struct {
	cache *cache.CodeCache
}

func NewCodeRepository(cache *cache.CodeCache) *CodeRepository {
	return &CodeRepository{cache: cache}
}

func (cr *CodeRepository) Set(ctx context.Context, biz, phone, code string) error {
	return cr.cache.Set(ctx, biz, phone, code)
}
