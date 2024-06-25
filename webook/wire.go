//go:build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/lcsin/goprojets/webook/internal/handler"
	"github.com/lcsin/goprojets/webook/internal/repository"
	"github.com/lcsin/goprojets/webook/internal/repository/cache"
	"github.com/lcsin/goprojets/webook/internal/repository/dao"
	"github.com/lcsin/goprojets/webook/internal/service"
	"github.com/lcsin/goprojets/webook/ioc"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		// 基础第三方服务
		ioc.InitDB, ioc.InitRedis,
		// dao
		dao.NewUserDAO, cache.NewUserCache, cache.NewCodeCache,
		// repository
		repository.NewUserRepository, repository.NewCodeRepository,
		// service
		ioc.InitSMSService,
		service.NewUserService, service.NewCodeService,
		// handler
		handler.NewUserHandler,
		// web server
		ioc.InitMiddlewares, ioc.InitWebServer,
	)
	return gin.Default()
}
