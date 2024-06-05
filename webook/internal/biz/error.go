package biz

import (
	"errors"
)

var (
	ErrDuplicateEmail      = errors.New("邮箱已注册")
	ErrInvalidUserOrPasswd = errors.New("用户不存在或密码不正确")
	ErrUserNotFound        = errors.New("用户不存在")
)
