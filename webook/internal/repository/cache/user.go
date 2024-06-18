package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/lcsin/goprojets/webook/internal/domain"
	"github.com/redis/go-redis/v9"
)

type IUserCache interface {
	Get(ctx context.Context, uid int64) (domain.User, error)
	Set(ctx context.Context, u domain.User) error
}

type UserCache struct {
	// 能用接口尽量用接口而不是直接用具体的实现，这样我们可以很灵活的使用这个接口
	// 例如我们可以自定义实现为Redis+本地缓存，当Redis崩了后启用本地缓存
	// 这里直接使用Redis的接口而不是Redis的Client
	client redis.Cmdable
	// 统一的过期时间
	expiration time.Duration
}

// NewUserCache 构造一个UserCache
// 面向接口编程和依赖注入的三个原则：
// 1. A 用到 B，B 一定是接口
// 2. A 用到 B，B 一定是A的字段
// 3. A 用到 B，A 一定不初始化B，而是外面注入
func NewUserCache(client redis.Cmdable) IUserCache {
	return &UserCache{client: client, expiration: time.Minute * 15}
}

func (uc *UserCache) key(uid int64) string {
	return fmt.Sprintf("user:info:%v", uid)
}

func (uc *UserCache) Get(ctx context.Context, uid int64) (domain.User, error) {
	result, err := uc.client.Get(ctx, uc.key(uid)).Result()
	if err != nil {
		return domain.User{}, err
	}

	var u domain.User
	if err = json.Unmarshal([]byte(result), &u); err != nil {
		return domain.User{}, err
	}
	return u, nil
}

func (uc *UserCache) Set(ctx context.Context, u domain.User) error {
	val, err := json.Marshal(&u)
	if err != nil {
		return err
	}

	return uc.client.Set(ctx, uc.key(u.UID), val, uc.expiration).Err()
}
