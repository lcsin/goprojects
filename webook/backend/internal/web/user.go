package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
}

func (u *UserHandler) RegisterRoutes(v1 *gin.RouterGroup) {
	ug := v1.Group("/users")
	ug.POST("/signup", u.Signup)
	ug.POST("/login", u.Login)
	ug.POST("/edit", u.Edit)
	ug.POST("/profile", u.Profile)
}

func (u *UserHandler) Signup(c *gin.Context) {
	c.String(http.StatusOK, "signup")
}

func (u *UserHandler) Login(c *gin.Context) {
}

func (u *UserHandler) Edit(c *gin.Context) {
}

func (u *UserHandler) Profile(c *gin.Context) {
}
