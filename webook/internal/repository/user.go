package repository

import (
	"context"
	"log"
	"time"

	"github.com/lcsin/goprojets/webook/internal/domain"
	"github.com/lcsin/goprojets/webook/internal/repository/cache"
	"github.com/lcsin/goprojets/webook/internal/repository/dao"
)

type UserRepository struct {
	dao   *dao.UserDAO     // 持久层
	cache *cache.UserCache // 缓存层
}

func NewUserRepository(dao *dao.UserDAO, cache *cache.UserCache) *UserRepository {
	return &UserRepository{dao: dao, cache: cache}
}

func (ur *UserRepository) Create(ctx context.Context, u domain.User) error {
	return ur.dao.Insert(ctx, dao.User{
		Email:  u.Email,
		Passwd: u.Passwd,
	})
}

func (ur *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	user, err := ur.dao.SelectByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}

	return domain.User{
		UID:        user.ID,
		Email:      user.Email,
		Passwd:     user.Passwd,
		Nickname:   user.Nickname,
		Profile:    user.Profile,
		CreateTime: time.UnixMilli(user.CreateTime),
		Birthday:   time.UnixMilli(user.Birthday),
	}, nil
}

func (ur *UserRepository) UpdateByID(ctx context.Context, u domain.User) error {
	return ur.dao.UpdateByID(ctx, dao.User{
		ID:       u.UID,
		Nickname: u.Nickname,
		Profile:  u.Profile,
		Birthday: u.Birthday.UnixMilli(),
	})
}

func (ur *UserRepository) FindByID(ctx context.Context, uid int64) (domain.User, error) {
	// 缓存中有直接返回
	cu, err := ur.cache.Get(ctx, uid)
	if err == nil {
		return cu, nil
	}

	// 从数据库中查询 - 这里需要考虑如果缓存崩了，会不会导致数据库也崩了
	du, err := ur.dao.SelectByID(ctx, uid)
	if err != nil {
		return domain.User{}, err
	}
	cu = domain.User{
		UID:        du.ID,
		Nickname:   du.Nickname,
		Email:      du.Email,
		Passwd:     du.Passwd,
		Profile:    du.Profile,
		Birthday:   time.UnixMilli(du.Birthday),
		CreateTime: time.UnixMilli(du.CreateTime),
	}

	// 写入到缓存
	if err = ur.cache.Set(ctx, cu); err != nil {
		log.Println("set user cache error: ", err)
	}

	return cu, nil
}
