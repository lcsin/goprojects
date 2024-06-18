package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var ErrCodeSendToMany = errors.New("验证码发送频繁")
var ErrCodeVerifyTooManyTimes = errors.New("验证码验证频繁")

type ICodeCache interface {
	Set(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz, phone, inputCode string) (bool, error)
}

type CodeCache struct {
	client redis.Cmdable
}

func NewCodeCache(client redis.Cmdable) ICodeCache {
	return &CodeCache{client: client}
}

//go:embed lua/set_code.lua
var luaSetCode string

//go:embed lua/verify_code.lua
var luaVerifyCode string

func (cc *CodeCache) Set(ctx context.Context, biz, phone, code string) error {
	res, err := cc.client.Eval(ctx, luaSetCode, []string{cc.key(biz, phone)}, code).Int()
	if err != nil {
		return err
	}

	switch res {
	case -1:
		return ErrCodeSendToMany
	case -2:
		return errors.New("系统错误")
	default:
		return nil
	}
}

func (cc *CodeCache) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	res, err := cc.client.Eval(ctx, luaVerifyCode, []string{cc.key(biz, phone)}, inputCode).Int()
	if err != nil {
		return false, err
	}

	switch res {
	case -1: // 验证太频繁，可能有人搞你
		return false, ErrCodeVerifyTooManyTimes
	case -2: // 验证失败，在允许范围内
		return false, nil
	default: // 验证通过
		return true, nil
	}
}

func (cc *CodeCache) key(biz, phone string) string {
	return fmt.Sprintf("phone_code:%v:%v", biz, phone)
}
