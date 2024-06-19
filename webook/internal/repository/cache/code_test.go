package cache

import (
	"context"
	"errors"
	"testing"

	redismocks "github.com/lcsin/goprojets/webook/internal/repository/cache/redismock"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCodeCache_Set(t *testing.T) {
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) redis.Cmdable

		// 输入
		ctx              context.Context
		biz, phone, code string
		// 输出
		wantErr error
	}{
		{
			name: "验证码设置成功",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redis.NewCmd(context.Background())
				cmd.SetVal(int64(0))
				cmdable := redismocks.NewMockCmdable(ctrl)
				cmdable.EXPECT().Eval(gomock.Any(), luaSetCode, []string{"phone_code:login:187xxx"}, []any{"123456"}).
					Return(cmd)
				return cmdable
			},
			ctx:     context.Background(),
			biz:     "login",
			phone:   "187xxx",
			code:    "123456",
			wantErr: nil,
		},
		{
			name: "redis错误",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redis.NewCmd(context.Background())
				cmd.SetErr(errors.New("redis error"))
				cmdable := redismocks.NewMockCmdable(ctrl)
				cmdable.EXPECT().Eval(gomock.Any(), luaSetCode, []string{"phone_code:login:187xxx"}, []any{"123456"}).
					Return(cmd)
				return cmdable
			},
			ctx:     context.Background(),
			biz:     "login",
			phone:   "187xxx",
			code:    "123456",
			wantErr: errors.New("redis error"),
		},
		{
			name: "验证码发送太频繁",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redis.NewCmd(context.Background())
				cmd.SetVal(int64(-1))
				cmdable := redismocks.NewMockCmdable(ctrl)
				cmdable.EXPECT().Eval(gomock.Any(), luaSetCode, []string{"phone_code:login:187xxx"}, []any{"123456"}).
					Return(cmd)
				return cmdable
			},
			ctx:     context.Background(),
			biz:     "login",
			phone:   "187xxx",
			code:    "123456",
			wantErr: ErrCodeSendToMany,
		},
		{
			name: "验证码发送太频繁",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redis.NewCmd(context.Background())
				cmd.SetVal(int64(-2))
				cmdable := redismocks.NewMockCmdable(ctrl)
				cmdable.EXPECT().Eval(gomock.Any(), luaSetCode, []string{"phone_code:login:187xxx"}, []any{"123456"}).
					Return(cmd)
				return cmdable
			},
			ctx:     context.Background(),
			biz:     "login",
			phone:   "187xxx",
			code:    "123456",
			wantErr: errors.New("系统错误"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			cc := NewCodeCache(tc.mock(ctrl))
			err := cc.Set(tc.ctx, tc.biz, tc.phone, tc.code)

			assert.Equal(t, tc.wantErr, err)
		})
	}
}
