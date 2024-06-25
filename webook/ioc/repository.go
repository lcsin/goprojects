package ioc

import (
	"github.com/lcsin/goprojets/webook/internal/repository/dao"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// InitDB 初始化MySQL
func InitDB() *gorm.DB {
	dns := viper.Get("mysql.dns").(string)
	db, err := gorm.Open(mysql.Open(dns))
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
	addr := viper.Get("redis.addr").(string)
	passwd := viper.Get("redis.passwd").(string)
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: passwd,
	})
}
