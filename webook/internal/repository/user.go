package repository

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/lcsin/goprojets/webook/internal/domain"
	"github.com/lcsin/goprojets/webook/internal/repository/cache"
	"github.com/lcsin/goprojets/webook/internal/repository/dao"
)

type IUserRepository interface {
	Create(ctx context.Context, u domain.User) error
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	UpdateByID(ctx context.Context, u domain.User) error
	FindByID(ctx context.Context, uid int64) (domain.User, error)
}

type UserRepository struct {
	dao   dao.IUserDAO     // 持久层
	cache cache.IUserCache // 缓存层
}

func NewUserRepository(dao dao.IUserDAO, cache cache.IUserCache) IUserRepository {
	return &UserRepository{dao: dao, cache: cache}
}

func (ur *UserRepository) Create(ctx context.Context, u domain.User) error {
	return ur.dao.Insert(ctx, ur.domain2Entity(u))
}

func (ur *UserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	user, err := ur.dao.SelectByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}

	return ur.entity2Domain(user), nil
}

func (ur *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	user, err := ur.dao.SelectByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}

	return ur.entity2Domain(user), nil
}

func (ur *UserRepository) UpdateByID(ctx context.Context, u domain.User) error {
	return ur.dao.UpdateByID(ctx, ur.domain2Entity(u))
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
	cu = ur.entity2Domain(du)

	// 写入到缓存
	if err = ur.cache.Set(ctx, cu); err != nil {
		log.Println("set user cache error: ", err)
	}

	return cu, nil
}

func (ur *UserRepository) entity2Domain(u dao.User) domain.User {
	return domain.User{
		UID:        u.ID,
		Nickname:   u.Nickname,
		Email:      u.Email.String,
		Phone:      u.Phone.String,
		Passwd:     u.Passwd,
		Profile:    u.Profile,
		Birthday:   time.UnixMilli(u.Birthday.Int64),
		CreateTime: time.UnixMilli(u.CreateTime),
	}
}

func (ur *UserRepository) domain2Entity(u domain.User) dao.User {
	return dao.User{
		ID: u.UID,
		Email: sql.NullString{
			String: u.Email,
			Valid:  u.Email != "",
		},
		Phone: sql.NullString{
			String: u.Phone,
			Valid:  u.Phone != "",
		},
		Passwd:   u.Passwd,
		Nickname: u.Nickname,
		Profile:  u.Profile,
		Birthday: sql.NullInt64{
			Int64: u.Birthday.UnixMilli(),
			Valid: u.Birthday.UnixMilli() > 0,
		},
		CreateTime: u.CreateTime.UnixMilli(),
	}
}
