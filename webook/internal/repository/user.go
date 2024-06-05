package repository

import (
	"context"
	"time"

	"github.com/lcsin/goprojets/webook/internal/domain"
	"github.com/lcsin/goprojets/webook/internal/repository/dao"
)

type UserRepository struct {
	dao *dao.UserDAO
}

func NewUserRepository(dao *dao.UserDAO) *UserRepository {
	return &UserRepository{dao: dao}
}

func (ur *UserRepository) Create(ctx context.Context, u domain.User) error {
	return ur.dao.Insert(ctx, dao.User{
		Email:  u.Email,
		Passwd: u.Passwd,
	})
}

func (ur *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := ur.dao.SelectByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return &domain.User{
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

func (ur *UserRepository) FindByID(ctx context.Context, ID int64) (*domain.User, error) {
	user, err := ur.dao.SelectByID(ctx, ID)
	if err != nil {
		return nil, err
	}
	return &domain.User{
		UID:        user.ID,
		Nickname:   user.Nickname,
		Email:      user.Email,
		Passwd:     user.Passwd,
		Profile:    user.Profile,
		Birthday:   time.UnixMilli(user.Birthday),
		CreateTime: time.UnixMilli(user.CreateTime),
	}, err
}
