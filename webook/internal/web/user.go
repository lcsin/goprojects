package web

import (
	regexp "github.com/dlclark/regexp2"
	"github.com/lcsin/goprojets/webook/internal/domain"
	"github.com/lcsin/goprojets/webook/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/lcsin/gopocket/util/ginx"
)

type UserHandler struct {
	srv            *service.UserService
	emailRexExp    *regexp.Regexp
	passwordRexExp *regexp.Regexp
}

func NewUserHandler(srv *service.UserService) *UserHandler {
	const (
		emailRegexPattern    = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
		passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
	)

	return &UserHandler{
		// 标准库的正则不支持复杂语法所以使用开源的正则库
		// 预编译正则，避免每次接口请求时都要编译
		emailRexExp:    regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordRexExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
		srv:            srv,
	}
}

func (u *UserHandler) RegisterRoutes(v1 *gin.RouterGroup) {
	ug := v1.Group("/users")
	ug.POST("/signup", u.Signup)
	ug.POST("/login", u.Login)
	ug.POST("/edit", u.Edit)
	ug.POST("/profile", u.Profile)
}

func (u *UserHandler) Signup(c *gin.Context) {
	type signupReq struct {
		Email         string `json:"email"`
		Passwd        string `json:"passwd"`
		ConfirmPasswd string `json:"confirmPasswd"`
	}
	var req signupReq
	if err := c.ShouldBind(&req); err != nil {
		ginx.Error(c, -400, "参数无效")
		return
	}
	// 参数校验
	isEmail, _ := u.emailRexExp.MatchString(req.Email)
	if !isEmail {
		ginx.Error(c, -400, "邮箱格式错误")
		return
	}
	isPasswd, _ := u.passwordRexExp.MatchString(req.Passwd)
	if !isPasswd {
		ginx.Error(c, -400, "密码必须包含数字、特殊字符，并且长度不能小于8位")
		return
	}
	if req.Passwd != req.ConfirmPasswd {
		ginx.Error(c, -400, "两次输入的密码不一致")
		return
	}
	// 业务注册
	if err := u.srv.Signup(c, domain.User{
		Email:  req.Email,
		Passwd: req.Passwd,
	}); err != nil {
		ginx.Error(c, -500, "系统错误")
		return
	}

	ginx.OK(c, nil)
}

func (u *UserHandler) Login(c *gin.Context) {
}

func (u *UserHandler) Edit(c *gin.Context) {
}

func (u *UserHandler) Profile(c *gin.Context) {
}
