package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/lcsin/goprojets/webook/internal/domain"
	"github.com/redis/go-redis/v9"
)

type UserCache struct {
	client     redis.Cmdable // 能用接口尽量用接口而不是直接用具体的实现
	expiration time.Duration
}

// NewUserCache 构造一个UserCache
// A 用到 B，B 一定是接口
// A 用到 B，B 一定是A的字段
// A 用到 B，A 一定不初始化B，而是外面注入
func NewUserCache(client redis.Cmdable, expiration time.Duration) *UserCache {
	return &UserCache{client: client, expiration: expiration}
}

func (uc *UserCache) key(uid int64) string {
	return fmt.Sprintf("user:info:%v", uid)
}

func (uc *UserCache) Get(ctx context.Context, uid int64) (domain.User, error) {
	result, err := uc.client.Get(ctx, uc.key(uid)).Result()
	if err != nil {
		return domain.User{}, nil
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
