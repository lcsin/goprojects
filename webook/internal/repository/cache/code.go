package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var ErrSendCodeToMany = errors.New("验证码发送频繁")

type CodeCache struct {
	client redis.Cmdable
}

func NewCodeCache(client redis.Cmdable) *CodeCache {
	return &CodeCache{client: client}
}

//go:embed lua/set_code.lua
var luaSetCode string

func (cc *CodeCache) Set(ctx context.Context, biz, phone, code string) error {
	res, err := cc.client.Eval(ctx, luaSetCode, []string{cc.key(biz, phone)}, code).Int()
	if err != nil {
		return err
	}

	switch res {
	case -1:
		return ErrSendCodeToMany
	case -2:
		return errors.New("系统错误")
	default:
		return nil
	}
}

func (cc *CodeCache) key(biz, phone string) string {
	return fmt.Sprintf("phone_code:%v:%v", biz, phone)
}
