package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/lcsin/gopocket/util/ginx"
	"github.com/lcsin/goprojets/webook/internal/integration/startup"
	"github.com/lcsin/goprojets/webook/internal/service"
	"github.com/lcsin/goprojets/webook/ioc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserHandler_e2e_SendLoginSMSCode(t *testing.T) {
	server := startup.InitWebServer()
	rdb := ioc.InitRedis()
	testCases := []struct {
		name string

		// 你要考虑准备数据。
		before func(t *testing.T)
		// 以及验证数据 数据库的数据对不对，你 Redis 的数据对不对
		after   func(t *testing.T)
		reqBody string

		wantCode int
		wantBody ginx.Response
	}{
		{
			name: "发送成功",
			before: func(t *testing.T) {
				// 不需要，也就是 Redis 什么数据也没有
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				// 你要清理数据
				// "phone_code:%s:%s"
				val, err := rdb.GetDel(ctx, "phone_code:login:15212345678").Result()
				cancel()
				assert.NoError(t, err)
				// 你的验证码是 6 位
				assert.True(t, len(val) == 6)
			},
			reqBody: `
{
	"phone": "15212345678"
}
`,
			wantCode: 200,
			wantBody: ginx.Response{
				Code:    0,
				State:   true,
				Message: "success",
				Data:    "code send success",
			},
		},
		{
			name: "发送太频繁",
			before: func(t *testing.T) {
				// 这个手机号码，已经有一个验证码了
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				_, err := rdb.Set(ctx, "phone_code:login:15212345678", "123456",
					time.Minute*9+time.Second*30).Result()
				cancel()
				assert.NoError(t, err)

			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				// 你要清理数据
				// "phone_code:%s:%s"
				val, err := rdb.GetDel(ctx, "phone_code:login:15212345678").Result()
				cancel()
				assert.NoError(t, err)
				// 你的验证码是 6 位,没有被覆盖，还是123456
				assert.Equal(t, "123456", val)
			},
			reqBody: `
		{
			"phone": "15212345678"
		}
		`,
			wantCode: 200,
			wantBody: ginx.Response{
				Code:    ginx.ErrBadRequest,
				State:   false,
				Message: service.ErrCodeSendTooMany.Error(),
				Data:    nil,
			},
		},
		{
			name: "系统错误",
			before: func(t *testing.T) {
				// 这个手机号码，已经有一个验证码了，但是没有过期时间
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				_, err := rdb.Set(ctx, "phone_code:login:15212345678", "123456", 0).Result()
				cancel()
				assert.NoError(t, err)

			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				// 你要清理数据
				// "phone_code:%s:%s"
				val, err := rdb.GetDel(ctx, "phone_code:login:15212345678").Result()
				cancel()
				assert.NoError(t, err)
				// 你的验证码是 6 位,没有被覆盖，还是123456
				assert.Equal(t, "123456", val)
			},
			reqBody: `
		{
			"phone": "15212345678"
		}
		`,
			wantCode: 200,
			wantBody: ginx.Response{
				Code:    ginx.ErrInternalServer,
				State:   false,
				Message: ginx.ErrInternalServer.String(),
				Data:    nil,
			},
		},
		{
			name: "手机号码为空",
			before: func(t *testing.T) {
			},
			after: func(t *testing.T) {
			},
			reqBody: `
		{
			"phone": ""
		}
		`,
			wantCode: 200,
			wantBody: ginx.Response{
				Code:    ginx.ErrBadRequest,
				State:   false,
				Message: "手机号码为空",
				Data:    nil,
			},
		},
		{
			name: "数据格式错误",
			before: func(t *testing.T) {
			},
			after: func(t *testing.T) {
			},
			reqBody: `
{
	"phone": ,
}
`,
			wantCode: 200,
			wantBody: ginx.Response{
				Code:    ginx.ErrBadRequest,
				State:   false,
				Message: ginx.ErrBadRequest.String(),
				Data:    nil,
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			req, err := http.NewRequest(http.MethodPost,
				"/api/v1/users/login/sms/code/send", bytes.NewBuffer([]byte(tc.reqBody)))
			require.NoError(t, err)
			// 数据是 JSON 格式
			req.Header.Set("Content-Type", "application/json")
			// 这里你就可以继续使用 req

			resp := httptest.NewRecorder()
			// 这就是 HTTP 请求进去 GIN 框架的入口。
			// 当你这样调用的时候，GIN 就会处理这个请求
			// 响应写回到 resp 里
			server.ServeHTTP(resp, req)

			assert.Equal(t, tc.wantCode, resp.Code)
			if resp.Code != 200 {
				return
			}
			var webRes ginx.Response
			err = json.NewDecoder(resp.Body).Decode(&webRes)
			require.NoError(t, err)
			assert.Equal(t, tc.wantBody, webRes)
			tc.after(t)
		})
	}
}
