package main

import (
	"github.com/lcsin/gopocket/util/httpx"
	"github.com/lcsin/goprojets/webook/internal/repository"
	"github.com/lcsin/goprojets/webook/internal/repository/dao"
	"github.com/lcsin/goprojets/webook/internal/service"
	"github.com/lcsin/goprojets/webook/internal/web"
)

func main() {
	db := repository.InitDB()
	ud := dao.NewUserDAO(db)
	ur := repository.NewUserRepository(ud)
	us := service.NewUserService(ur)
	r := web.RegisterRoutes(us)

	httpx.Graceful(r, ":8080")
}
