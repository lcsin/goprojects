package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lcsin/gopocket/util/httpx"
)

func main() {
	//db := repository.InitDB()
	//ud := dao.NewUserDAO(db)
	//ur := repository.NewUserRepository(ud)
	//us := service.NewUserService(ur)
	//r := web.RegisterRoutes(us)

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	httpx.Graceful(r, ":8080")
}
