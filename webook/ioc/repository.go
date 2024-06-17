package ioc

import (
	"github.com/lcsin/goprojets/webook/config"
	"github.com/lcsin/goprojets/webook/internal/repository/dao"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// InitDB 初始化MySQL
func InitDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open(config.Cfg.DB.DNS))
	if err != nil {
		panic(err)
	}

	if err = db.AutoMigrate(dao.User{}); err != nil {
		panic(err)
	}

	return db
}

// InitRedis 初始化Redis
func InitRedis() redis.Cmdable {
	return redis.NewClient(&redis.Options{
		Addr:     config.Cfg.Redis.Addr,
		Password: config.Cfg.Redis.Passwd,
	})
}
