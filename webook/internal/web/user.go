package web

import (
	"errors"
	"fmt"
	"time"

	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lcsin/gopocket/util/ginx"
	"github.com/lcsin/goprojets/webook/internal/biz"
	"github.com/lcsin/goprojets/webook/internal/domain"
	"github.com/lcsin/goprojets/webook/internal/service"
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

// Signup 用户注册
func (u *UserHandler) Signup(c *gin.Context) {
	type signupReq struct {
		Email         string `json:"email"`
		Passwd        string `json:"passwd"`
		ConfirmPasswd string `json:"confirmPasswd"`
	}
	var req signupReq
	if err := c.ShouldBind(&req); err != nil {
		ginx.ResponseError(c, ginx.ErrBadRequest)
		return
	}
	// 参数校验
	isEmail, _ := u.emailRexExp.MatchString(req.Email)
	if !isEmail {
		ginx.ResponseErrorMessage(c, ginx.ErrBadRequest, "邮箱格式错误")
		return
	}
	isPasswd, _ := u.passwordRexExp.MatchString(req.Passwd)
	if !isPasswd {
		ginx.ResponseErrorMessage(c, ginx.ErrBadRequest, "密码必须包含数字、特殊字符，并且长度不能小于8位")
		return
	}
	if req.Passwd != req.ConfirmPasswd {
		ginx.ResponseErrorMessage(c, ginx.ErrBadRequest, "两次输入的密码不一致")
		return
	}
	// 用户注册
	err := u.srv.Signup(c, domain.User{
		Email:  req.Email,
		Passwd: req.Passwd,
	})

	switch err {
	case nil:
		ginx.ResponseOK(c, nil)
	case biz.ErrDuplicateEmail:
		ginx.ResponseErrorMessage(c, ginx.ErrBadRequest, err.Error())
	default:
		ginx.ResponseError(c, ginx.ErrInternalServer)
	}
}

// Login 用户登录
func (u *UserHandler) Login(c *gin.Context) {
	type loginReq struct {
		Email  string `json:"'email'"`
		Passwd string `json:"passwd"`
	}
	var req loginReq
	if err := c.ShouldBind(&req); err != nil {
		ginx.ResponseError(c, ginx.ErrBadRequest)
		return
	}
	// 用户登录
	user, err := u.srv.Login(c, domain.User{
		Email:  req.Email,
		Passwd: req.Passwd,
	})
	if err != nil {
		if errors.Is(err, biz.ErrInvalidUserOrPasswd) {
			ginx.ResponseErrorMessage(c, ginx.ErrBadRequest, err.Error())
			return
		}
		ginx.ResponseError(c, ginx.ErrInternalServer)
		return
	}

	// 生成jwt
	claims := biz.UserClaims{UID: user.UID}
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7))
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(biz.JwtKey))
	if err != nil {
		ginx.ResponseError(c, ginx.ErrInternalServer)
		return
	}
	// 设置jwt token
	c.Header("x-jwt-token", fmt.Sprintf("Bearer %v", token))
	ginx.ResponseOK(c, user)
}

func (u *UserHandler) Edit(c *gin.Context) {
	type editReq struct {
		UID      int64  `json:"uid"`
		Nickname string `json:"nickname"`
		Email    string `json:"email"`
		Passwd   string `json:"passwd"`
		Profile  string `json:"profile"`
		Birthday string `json:"birthday"`
	}
	var req editReq
	if err := c.ShouldBind(&req); err != nil {
		ginx.ResponseError(c, ginx.ErrBadRequest)
		return
	}

	birthday, err := time.Parse(time.DateOnly, req.Birthday)
	if err != nil {
		ginx.ResponseError(c, ginx.ErrBadRequest)
		return
	}

	err = u.srv.Edit(c, domain.User{
		UID:      req.UID,
		Nickname: req.Nickname,
		Email:    req.Email,
		Passwd:   req.Passwd,
		Profile:  req.Profile,
		Birthday: birthday,
	})
	if err != nil {
		if errors.Is(err, biz.ErrUserNotFound) {
			ginx.ResponseErrorMessage(c, ginx.ErrBadRequest, err.Error())
			return
		}
		ginx.ResponseError(c, ginx.ErrInternalServer)
		return
	}

	ginx.ResponseOK(c, nil)
}

func (u *UserHandler) Profile(c *gin.Context) {
	uid, ok := c.Get("uid")
	if !ok {
		ginx.ResponseError(c, ginx.ErrNotFound)
		return
	}

	user, err := u.srv.Profile(c, uid.(int64))
	if err != nil {
		if errors.Is(err, biz.ErrUserNotFound) {
			ginx.ResponseErrorMessage(c, ginx.ErrNotFound, err.Error())
			return
		}
		ginx.ResponseError(c, ginx.ErrInternalServer)
		return
	}

	ginx.ResponseOK(c, user)
}
