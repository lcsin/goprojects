package service

import (
	"context"
	"errors"

	"github.com/lcsin/goprojets/webook/internal/biz"
	"github.com/lcsin/goprojets/webook/internal/domain"
	"github.com/lcsin/goprojets/webook/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (us *UserService) Signup(ctx context.Context, u domain.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Passwd), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Passwd = string(hash)

	return us.repo.Create(ctx, u)
}

func (us *UserService) Login(ctx context.Context, u domain.User) (*domain.User, error) {
	user, err := us.repo.FindByEmail(ctx, u.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, biz.ErrInvalidUserOrPasswd
		}
		return nil, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Passwd), []byte(u.Passwd)); err != nil {
		return nil, biz.ErrInvalidUserOrPasswd
	}

	user.Passwd = ""
	return user, nil
}

func (us *UserService) Edit(ctx context.Context, u domain.User) error {
	err := us.repo.UpdateByID(ctx, u)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return biz.ErrUserNotFound
		}
	}
	return err
}

func (us *UserService) Profile(ctx context.Context, ID int64) (*domain.User, error) {
	user, err := us.repo.FindByID(ctx, ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, biz.ErrUserNotFound
		}
		return nil, err
	}
	user.Passwd = ""
	return user, nil
}
