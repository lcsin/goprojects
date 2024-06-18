package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/lcsin/goprojets/webook/internal/domain"
	"github.com/lcsin/goprojets/webook/internal/service"
	svcmocks "github.com/lcsin/goprojets/webook/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestUserHandler_Signup(t *testing.T) {
	testCases := []struct {
		name    string
		mock    func(ctrl *gomock.Controller) service.IUserService
		reqBody string
		bizCode float64
		bizMsg  string
	}{
		{
			name: "注册成功",
			mock: func(ctrl *gomock.Controller) service.IUserService {
				usersvc := svcmocks.NewMockIUserService(ctrl)
				usersvc.EXPECT().Signup(gomock.Any(), domain.User{
					Email:  "123@qq.com",
					Passwd: "hello@world123",
				}).Return(nil)

				return usersvc
			},
			reqBody: `
{
	"email": "123@qq.com",
	"passwd": "hello@world123",
	"confirmPasswd": "hello@world123"
}
`,
			bizCode: 0,
			bizMsg:  "success",
		},
		{
			name: "邮箱格式错误",
			mock: func(ctrl *gomock.Controller) service.IUserService {
				usersvc := svcmocks.NewMockIUserService(ctrl)
				return usersvc
			},
			reqBody: `
{
	"email": "123@q123",
	"passwd": "hello@world123",
	"confirmPasswd": "hello@world123"
}
`,
			bizCode: -400,
			bizMsg:  "邮箱格式错误",
		},
		{
			name: "密码必须包含数字、特殊字符，并且长度不能小于8位",
			mock: func(ctrl *gomock.Controller) service.IUserService {
				usersvc := svcmocks.NewMockIUserService(ctrl)
				return usersvc
			},
			reqBody: `
{
	"email": "123@qq.com",
	"passwd": "hello123",
	"confirmPasswd": "hello123"
}
`,
			bizCode: -400,
			bizMsg:  "密码必须包含数字、特殊字符，并且长度不能小于8位",
		},
		{
			name: "两次输入的密码不一致",
			mock: func(ctrl *gomock.Controller) service.IUserService {
				usersvc := svcmocks.NewMockIUserService(ctrl)
				return usersvc
			},
			reqBody: `
{
	"email": "123@qq.com",
	"passwd": "hello@world123",
	"confirmPasswd": "hello@world"
}
`,
			bizCode: -400,
			bizMsg:  "两次输入的密码不一致",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			server := gin.Default()
			h := NewUserHandler(tc.mock(ctrl), nil)
			h.RegisterRoutes(server.Group("/api/v1"))

			req, err := http.NewRequest(http.MethodPost, "/api/v1/users/signup", bytes.NewBuffer([]byte(tc.reqBody)))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			server.ServeHTTP(resp, req)

			var rep map[string]interface{}
			err = json.Unmarshal(resp.Body.Bytes(), &rep)
			require.NoError(t, err)

			assert.Equal(t, tc.bizCode, rep["code"])
			assert.Equal(t, tc.bizMsg, rep["message"])
		})
	}
}
