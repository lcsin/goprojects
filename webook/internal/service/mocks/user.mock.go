// Code generated by MockGen. DO NOT EDIT.
// Source: internal/service/user.go
//
// Generated by this command:
//
//	mockgen -source=internal/service/user.go -package=svcmocks -destination=internal/service/mocks/user.mock.go
//

// Package svcmocks is a generated GoMock package.
package svcmocks

import (
	context "context"
	reflect "reflect"

	domain "github.com/lcsin/goprojets/webook/internal/domain"
	gomock "go.uber.org/mock/gomock"
)

// MockIUserService is a mock of IUserService interface.
type MockIUserService struct {
	ctrl     *gomock.Controller
	recorder *MockIUserServiceMockRecorder
}

// MockIUserServiceMockRecorder is the mock recorder for MockIUserService.
type MockIUserServiceMockRecorder struct {
	mock *MockIUserService
}

// NewMockIUserService creates a new mock instance.
func NewMockIUserService(ctrl *gomock.Controller) *MockIUserService {
	mock := &MockIUserService{ctrl: ctrl}
	mock.recorder = &MockIUserServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIUserService) EXPECT() *MockIUserServiceMockRecorder {
	return m.recorder
}

// Edit mocks base method.
func (m *MockIUserService) Edit(ctx context.Context, u domain.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Edit", ctx, u)
	ret0, _ := ret[0].(error)
	return ret0
}

// Edit indicates an expected call of Edit.
func (mr *MockIUserServiceMockRecorder) Edit(ctx, u any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Edit", reflect.TypeOf((*MockIUserService)(nil).Edit), ctx, u)
}

// FindOrCreate mocks base method.
func (m *MockIUserService) FindOrCreate(ctx context.Context, phone string) (domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindOrCreate", ctx, phone)
	ret0, _ := ret[0].(domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindOrCreate indicates an expected call of FindOrCreate.
func (mr *MockIUserServiceMockRecorder) FindOrCreate(ctx, phone any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindOrCreate", reflect.TypeOf((*MockIUserService)(nil).FindOrCreate), ctx, phone)
}

// Login mocks base method.
func (m *MockIUserService) Login(ctx context.Context, email, passwd string) (domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", ctx, email, passwd)
	ret0, _ := ret[0].(domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Login indicates an expected call of Login.
func (mr *MockIUserServiceMockRecorder) Login(ctx, email, passwd any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockIUserService)(nil).Login), ctx, email, passwd)
}

// Profile mocks base method.
func (m *MockIUserService) Profile(ctx context.Context, uid int64) (domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Profile", ctx, uid)
	ret0, _ := ret[0].(domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Profile indicates an expected call of Profile.
func (mr *MockIUserServiceMockRecorder) Profile(ctx, uid any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Profile", reflect.TypeOf((*MockIUserService)(nil).Profile), ctx, uid)
}

// Signup mocks base method.
func (m *MockIUserService) Signup(ctx context.Context, u domain.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Signup", ctx, u)
	ret0, _ := ret[0].(error)
	return ret0
}

// Signup indicates an expected call of Signup.
func (mr *MockIUserServiceMockRecorder) Signup(ctx, u any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Signup", reflect.TypeOf((*MockIUserService)(nil).Signup), ctx, u)
}
