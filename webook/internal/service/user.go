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

func (us *UserService) Login(ctx context.Context, u domain.User) (domain.User, error) {
	user, err := us.repo.FindByEmail(ctx, u.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, biz.ErrInvalidUserOrPasswd
		}
		return domain.User{}, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Passwd), []byte(u.Passwd)); err != nil {
		return domain.User{}, biz.ErrInvalidUserOrPasswd
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

func (us *UserService) Profile(ctx context.Context, uid int64) (domain.User, error) {
	user, err := us.repo.FindByID(ctx, uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, biz.ErrUserNotFound
		}
		return domain.User{}, err
	}
	user.Passwd = ""
	return user, nil
}

func (us *UserService) FindOrCreate(ctx context.Context, phone string) (domain.User, error) {
	// 快路径: 只经过一个查询
	// 从业务角度考虑，这段代码可以删掉，不影响功能
	// 从现实角度考虑，加上这段代码可以有效的提高性能
	// 因为大部分用户是已经注册过的，这样他们通过手机号登录时，只需要一个查询就可以
	user, err := us.repo.FindByPhone(ctx, phone)
	if err != gorm.ErrRecordNotFound {
		return user, err
	}

	// 慢路径：需要走两个查询
	err = us.repo.Create(ctx, domain.User{Phone: phone})
	if err != nil && err != biz.ErrDuplicate {
		return user, err
	}

	// 多库情况下会有主从延迟的问题
	return us.repo.FindByPhone(ctx, phone)
}
