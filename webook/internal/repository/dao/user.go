package dao

import (
	"context"
	"database/sql"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/lcsin/goprojets/webook/internal/biz"
	"gorm.io/gorm"
)

type IUserDAO interface {
	Insert(ctx context.Context, u User) error
	SelectByPhone(ctx context.Context, phone string) (User, error)
	SelectByEmail(ctx context.Context, email string) (User, error)
	SelectByID(ctx context.Context, ID int64) (User, error)
	UpdateByID(ctx context.Context, u User) error
}

type UserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) IUserDAO {
	return &UserDAO{db: db}
}

type User struct {
	ID int64 `gorm:"primaryKey,autoIncrement"`
	// 唯一索引允许有多个空值
	Email    sql.NullString `gorm:"unique"`
	Phone    sql.NullString `gorm:"unique"`
	Passwd   string
	Nickname string
	Profile  string
	Birthday sql.NullInt64

	// 毫秒数
	CreateTime int64
	UpdateTime int64
}

func (dao *UserDAO) Insert(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.CreateTime = now
	u.UpdateTime = now
	err := dao.db.WithContext(ctx).Create(&u).Error
	if me, ok := err.(*mysql.MySQLError); ok {
		const duplicateErr uint16 = 1062
		if me.Number == duplicateErr {
			// 邮箱已被注册
			return biz.ErrDuplicate
		}
	}

	return err
}

func (dao *UserDAO) SelectByPhone(ctx context.Context, phone string) (User, error) {
	var user User
	if err := dao.db.WithContext(ctx).Where("phone = ?", phone).First(&user).Error; err != nil {
		return User{}, err
	}

	return user, nil
}

func (dao *UserDAO) SelectByEmail(ctx context.Context, email string) (User, error) {
	var user User
	if err := dao.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return User{}, err
	}
	return user, nil
}

func (dao *UserDAO) SelectByID(ctx context.Context, ID int64) (User, error) {
	var user User
	if err := dao.db.WithContext(ctx).Where("id = ?", ID).First(&user).Error; err != nil {
		return User{}, err
	}
	return user, nil
}

func (dao *UserDAO) UpdateByID(ctx context.Context, u User) error {
	if _, err := dao.SelectByID(ctx, u.ID); err != nil {
		return err
	}

	u.UpdateTime = time.Now().UnixMilli()
	return dao.db.WithContext(ctx).Where("id = ?", u.ID).Updates(&u).Error
}
