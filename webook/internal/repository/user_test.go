package repository

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/lcsin/goprojets/webook/internal/domain"
	"github.com/lcsin/goprojets/webook/internal/repository/cache"
	cachemocks "github.com/lcsin/goprojets/webook/internal/repository/cache/mocks"
	"github.com/lcsin/goprojets/webook/internal/repository/dao"
	daomocks "github.com/lcsin/goprojets/webook/internal/repository/dao/mocks"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUserRepository_FindByID(t *testing.T) {
	now := time.Now()
	now = time.UnixMilli(now.UnixMilli())
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) (dao.IUserDAO, cache.IUserCache)

		// 输入
		ctx context.Context
		id  int64
		// 输出
		wantUser domain.User
		wanErr   error
	}{
		{
			name: "查找缓存成功",
			mock: func(ctrl *gomock.Controller) (dao.IUserDAO, cache.IUserCache) {
				ud := daomocks.NewMockIUserDAO(ctrl)
				uc := cachemocks.NewMockIUserCache(ctrl)
				uc.EXPECT().Get(gomock.Any(), int64(1)).
					Return(domain.User{
						UID:      1,
						Nickname: "test",
						Email:    "1847@qq.com",
					}, nil)
				return ud, uc
			},
			ctx:    context.Background(),
			id:     1,
			wanErr: nil,
			wantUser: domain.User{
				UID:      1,
				Nickname: "test",
				Email:    "1847@qq.com",
			},
		},
		{
			name: "缓存未命中，数据库查找成功，并成功写入缓存",
			mock: func(ctrl *gomock.Controller) (dao.IUserDAO, cache.IUserCache) {
				uc := cachemocks.NewMockIUserCache(ctrl)
				uc.EXPECT().Get(gomock.Any(), int64(1)).
					Return(domain.User{}, redis.Nil)

				ud := daomocks.NewMockIUserDAO(ctrl)
				ud.EXPECT().SelectByID(gomock.Any(), int64(1)).
					Return(dao.User{
						ID:       1,
						Nickname: "test",
						Email: sql.NullString{
							String: "1847@qq.com",
							Valid:  true,
						},
						CreateTime: now.UnixMilli(),
						Birthday: sql.NullInt64{
							Int64: now.UnixMilli(),
							Valid: true,
						},
					}, nil)

				uc.EXPECT().Set(gomock.Any(), domain.User{
					UID:        1,
					Nickname:   "test",
					Email:      "1847@qq.com",
					CreateTime: now,
					Birthday:   now,
				}).Return(nil)

				return ud, uc
			},
			ctx:    context.Background(),
			id:     1,
			wanErr: nil,
			wantUser: domain.User{
				UID:        1,
				Nickname:   "test",
				Email:      "1847@qq.com",
				CreateTime: now,
				Birthday:   now,
			},
		},
		{
			name: "缓存未命中，数据库查询失败",
			mock: func(ctrl *gomock.Controller) (dao.IUserDAO, cache.IUserCache) {
				uc := cachemocks.NewMockIUserCache(ctrl)
				uc.EXPECT().Get(gomock.Any(), int64(1)).
					Return(domain.User{}, redis.Nil)

				ud := daomocks.NewMockIUserDAO(ctrl)
				ud.EXPECT().SelectByID(gomock.Any(), int64(1)).
					Return(dao.User{}, errors.New("db error"))
				return ud, uc
			},
			ctx:      context.Background(),
			id:       1,
			wanErr:   errors.New("db error"),
			wantUser: domain.User{},
		},
		{
			name: "缓存未命中，数据库查找成功，但写入缓存失败",
			mock: func(ctrl *gomock.Controller) (dao.IUserDAO, cache.IUserCache) {
				uc := cachemocks.NewMockIUserCache(ctrl)
				uc.EXPECT().Get(gomock.Any(), int64(1)).
					Return(domain.User{}, redis.Nil)

				ud := daomocks.NewMockIUserDAO(ctrl)
				ud.EXPECT().SelectByID(gomock.Any(), int64(1)).
					Return(dao.User{
						ID:       1,
						Nickname: "test",
						Email: sql.NullString{
							String: "1847@qq.com",
							Valid:  true,
						},
						CreateTime: now.UnixMilli(),
						Birthday: sql.NullInt64{
							Int64: now.UnixMilli(),
							Valid: true,
						},
					}, nil)

				uc.EXPECT().Set(gomock.Any(), domain.User{
					UID:        1,
					Nickname:   "test",
					Email:      "1847@qq.com",
					CreateTime: now,
					Birthday:   now,
				}).Return(errors.New("set cache error"))

				return ud, uc
			},
			ctx:    context.Background(),
			id:     1,
			wanErr: nil,
			wantUser: domain.User{
				UID:        1,
				Nickname:   "test",
				Email:      "1847@qq.com",
				CreateTime: now,
				Birthday:   now,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ud, uc := tc.mock(ctrl)
			repo := NewUserRepository(ud, uc)
			user, err := repo.FindByID(tc.ctx, tc.id)

			assert.Equal(t, tc.wanErr, err)
			assert.Equal(t, tc.wantUser, user)
		})
	}
}
