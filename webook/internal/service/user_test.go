package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/lcsin/goprojets/webook/internal/biz"
	"github.com/lcsin/goprojets/webook/internal/domain"
	"github.com/lcsin/goprojets/webook/internal/repository"
	repomocks "github.com/lcsin/goprojets/webook/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func TestUserService_Login(t *testing.T) {
	now := time.Now()
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) repository.IUserRepository
		// 输入
		ctx           context.Context
		email, passwd string
		// 输出
		wantUser domain.User
		wantErr  error
	}{
		{
			name: "登录成功",
			mock: func(ctrl *gomock.Controller) repository.IUserRepository {
				repo := repomocks.NewMockIUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "1847@qq.com").
					Return(domain.User{
						UID:        1,
						Nickname:   "test",
						Passwd:     "$2a$10$kdry4.pWXckh8frseTAqMuo9JuYZ4qFRVIMvoxiDrzaiHEt3Eh8gW",
						Email:      "1847@qq.com",
						CreateTime: now,
					}, nil)
				return repo
			},
			email:  "1847@qq.com",
			passwd: "hello@world123",
			wantUser: domain.User{
				UID:        1,
				Nickname:   "test",
				Passwd:     "$2a$10$kdry4.pWXckh8frseTAqMuo9JuYZ4qFRVIMvoxiDrzaiHEt3Eh8gW",
				Email:      "1847@qq.com",
				CreateTime: now,
			},
			wantErr: nil,
		},
		{
			name: "用户不存在",
			mock: func(ctrl *gomock.Controller) repository.IUserRepository {
				repo := repomocks.NewMockIUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "1847@qq.com").
					Return(domain.User{}, gorm.ErrRecordNotFound)
				return repo
			},
			email:    "1847@qq.com",
			passwd:   "hello@world123",
			wantUser: domain.User{},
			wantErr:  biz.ErrInvalidUserOrPasswd,
		},
		{
			name: "密码不正确",
			mock: func(ctrl *gomock.Controller) repository.IUserRepository {
				repo := repomocks.NewMockIUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "1847@qq.com").
					Return(domain.User{}, nil)
				return repo
			},
			email:    "1847@qq.com",
			passwd:   "hello@world123",
			wantUser: domain.User{},
			wantErr:  biz.ErrInvalidUserOrPasswd,
		},
		{
			name: "系统错误",
			mock: func(ctrl *gomock.Controller) repository.IUserRepository {
				repo := repomocks.NewMockIUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "1847@qq.com").
					Return(domain.User{}, errors.New("系统错误"))
				return repo
			},
			email:    "1847@qq.com",
			passwd:   "hello@world123",
			wantUser: domain.User{},
			wantErr:  errors.New("系统错误"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			svc := NewUserService(tc.mock(ctrl))
			user, err := svc.Login(tc.ctx, tc.email, tc.passwd)

			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUser, user)
		})
	}
}

func TestEncrypt(t *testing.T) {
	passwd := "hello@world123"
	res, _ := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost)
	t.Log(string(res))
}
