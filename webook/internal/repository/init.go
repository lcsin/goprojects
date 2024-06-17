package repository

import (
	"github.com/lcsin/goprojets/webook/internal/config"
	"github.com/lcsin/goprojets/webook/internal/repository/dao"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// InitDB 初始化数据库
func InitDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open(config.Config.DB.DNS))
	if err != nil {
		panic(err)
	}

	if err = db.AutoMigrate(dao.User{}); err != nil {
		panic(err)
	}

	return db
}

func InitRedis() redis.Cmdable {
	return redis.NewClient(&redis.Options{
		Addr:     config.Config.Redis.Addr,
		Password: config.Config.Redis.Passwd,
	})
}
