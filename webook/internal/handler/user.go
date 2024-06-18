package handler

import (
	"errors"
	"fmt"
	"time"

	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lcsin/gopocket/util/ginx"
	"github.com/lcsin/goprojets/webook/config"
	"github.com/lcsin/goprojets/webook/internal/biz"
	"github.com/lcsin/goprojets/webook/internal/domain"
	"github.com/lcsin/goprojets/webook/internal/service"
)

const (
	bizKey = "login"
)

type UserHandler struct {
	srv            service.IUserService
	codeSrv        service.ICodeService
	emailRexExp    *regexp.Regexp
	passwordRexExp *regexp.Regexp
}

type UserClaims struct {
	jwt.RegisteredClaims
	UID       int64  `json:"uid"`
	UserAgent string `json:"userAgent"`
}

func NewUserHandler(srv service.IUserService, codeSrv service.ICodeService) *UserHandler {
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
		codeSrv:        codeSrv,
	}
}

func (u *UserHandler) RegisterRoutes(v1 *gin.RouterGroup) {
	ug := v1.Group("/users")
	ug.POST("/signup", u.Signup)
	ug.POST("/login", u.Login)
	ug.POST("/edit", u.Edit)
	ug.POST("/profile", u.Profile)
	ug.POST("/login/sms/code/send", u.SendLoginSMSCode)
	ug.POST("/login/sms", u.LoginSMS)
}

func (u *UserHandler) LoginSMS(c *gin.Context) {
	type LoginSMS struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}
	var req LoginSMS
	if err := c.ShouldBind(&req); err != nil {
		ginx.ResponseError(c, ginx.ErrBadRequest)
		return
	}

	verify, err := u.codeSrv.Verify(c, bizKey, req.Phone, req.Code)
	if err != nil {
		ginx.ResponseErrorMessage(c, ginx.ErrBadRequest, err.Error())
		return
	}
	if !verify {
		ginx.ResponseErrorMessage(c, ginx.ErrBadRequest, "验证码错误")
		return
	}

	user, err := u.srv.FindOrCreate(c, req.Phone)
	if err != nil {
		ginx.ResponseError(c, ginx.ErrInternalServer)
		return
	}
	if err = u.setJWTToken(c, user); err != nil {
		ginx.ResponseError(c, ginx.ErrInternalServer)
		return
	}

	ginx.ResponseOK(c, "login success")
}

func (u *UserHandler) SendLoginSMSCode(c *gin.Context) {
	type SendLoginSMSCodeReq struct {
		Phone string `json:"phone"`
	}
	var req SendLoginSMSCodeReq
	if err := c.ShouldBind(&req); err != nil {
		ginx.ResponseError(c, ginx.ErrBadRequest)
		return
	}

	err := u.codeSrv.Send(c, bizKey, req.Phone)
	switch err {
	case nil:
		ginx.ResponseOK(c, "code send success")
	case service.ErrCodeSendTooMany:
		ginx.ResponseErrorMessage(c, ginx.ErrBadRequest, err.Error())
	default:
		ginx.ResponseError(c, ginx.ErrInternalServer)
	}
}

// Signup 用户注册
func (u *UserHandler) Signup(c *gin.Context) {
	type SignupReq struct {
		Email         string `json:"email"`
		Passwd        string `json:"passwd"`
		ConfirmPasswd string `json:"confirmPasswd"`
	}
	var req SignupReq
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
	case biz.ErrDuplicate:
		ginx.ResponseErrorMessage(c, ginx.ErrBadRequest, err.Error())
	default:
		ginx.ResponseError(c, ginx.ErrInternalServer)
	}
}

// Login 用户登录
func (u *UserHandler) Login(c *gin.Context) {
	type LoginReq struct {
		Email  string `json:"'email'"`
		Passwd string `json:"passwd"`
	}
	var req LoginReq
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
	if err = u.setJWTToken(c, user); err != nil {
		return
	}
	ginx.ResponseOK(c, user)
}

func (u *UserHandler) setJWTToken(c *gin.Context, user domain.User) error {
	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		},
		UID:       user.UID,
		UserAgent: c.GetHeader("User-Agent"),
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(config.Cfg.JWTKey))
	if err != nil {
		ginx.ResponseError(c, ginx.ErrInternalServer)
		return err
	}
	// 设置jwt token
	c.Header("x-jwt-token", fmt.Sprintf("Bearer %v", token))
	return nil
}

func (u *UserHandler) Edit(c *gin.Context) {
	type EditReq struct {
		UID      int64  `json:"uid"`
		Nickname string `json:"nickname"`
		Email    string `json:"email"`
		Passwd   string `json:"passwd"`
		Profile  string `json:"profile"`
		Birthday string `json:"birthday"`
	}
	var req EditReq
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
	type Profile struct {
		UID          int64  `json:"uid"`
		Email        string `json:"email"`
		Phone        string `json:"phone"`
		Nickname     string `json:"nickname"`
		AboutMe      string `json:"aboutMe"`
		Birthday     int64  `json:"birthday"`
		RegisterTime int64  `json:"registerTime"`
	}

	uid := c.MustGet("uid").(int64)
	user, err := u.srv.Profile(c, uid)
	if err != nil {
		if errors.Is(err, biz.ErrUserNotFound) {
			ginx.ResponseErrorMessage(c, ginx.ErrNotFound, err.Error())
			return
		}
		ginx.ResponseError(c, ginx.ErrInternalServer)
		return
	}

	profile := &Profile{
		UID:          user.UID,
		Email:        user.Email,
		Phone:        user.Phone,
		Nickname:     user.Nickname,
		AboutMe:      user.Profile,
		Birthday:     user.Birthday.UnixMilli(),
		RegisterTime: user.CreateTime.UnixMilli(),
	}

	ginx.ResponseOK(c, profile)
}
