package main

import (
	"github.com/lcsin/gopocket/util/httpx"
	"github.com/lcsin/goprojets/webook/internal/repository"
	"github.com/lcsin/goprojets/webook/internal/repository/cache"
	"github.com/lcsin/goprojets/webook/internal/repository/dao"
	"github.com/lcsin/goprojets/webook/internal/service"
	"github.com/lcsin/goprojets/webook/internal/service/sms/local"
	"github.com/lcsin/goprojets/webook/internal/web"
)

func main() {
	db := repository.InitDB()
	userDao := dao.NewUserDAO(db)

	rdb := repository.InitRedis()
	userCache := cache.NewUserCache(rdb)
	codeCache := cache.NewCodeCache(rdb)

	ur := repository.NewUserRepository(userDao, userCache)
	cr := repository.NewCodeRepository(codeCache)

	us := service.NewUserService(ur)
	sms := local.NewService()
	cs := service.NewCodeService(cr, sms)

	r := web.RegisterRoutes(us, cs)

	httpx.Graceful(r, ":8080")
}
