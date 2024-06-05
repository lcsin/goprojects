package dao

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type UserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{db: db}
}

type User struct {
	Id     int64  `gorm:"primaryKey,autoIncrement"`
	Email  string `gorm:"unique"`
	Passwd string

	// 毫秒数
	CreateTime int64
	UpdateTime int64
}

func (dao *UserDAO) Insert(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.CreateTime = now
	u.UpdateTime = now
	return dao.db.WithContext(ctx).Create(&u).Error
}
