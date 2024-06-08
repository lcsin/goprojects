package repository

import (
	"github.com/lcsin/goprojets/webook/internal/repository/dao"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// InitDB 初始化数据库
func InitDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13306)/webook"))
	if err != nil {
		panic(err)
	}

	if err = db.AutoMigrate(dao.User{}); err != nil {
		panic(err)
	}

	return db
}
