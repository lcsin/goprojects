package biz

import (
	"errors"
)

var (
	ErrDuplicate           = errors.New("用户已存在")
	ErrUserNotFound        = errors.New("用户不存在")
	ErrInvalidUserOrPasswd = errors.New("用户不存在或密码不正确")
)
