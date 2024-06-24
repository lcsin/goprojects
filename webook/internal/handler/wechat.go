package handler

import "github.com/gin-gonic/gin"

type WechatHandler struct {
}

func (w *WechatHandler) RegisterRoutes(v1 *gin.RouterGroup) {
	wechat := v1.Group("/wechat")
	wechat.GET("/oauth2/authurl", w.AuthURL)
	wechat.Any("/callback", w.Callback)
}

// AuthURL 构造跳到微信的URL
func (w *WechatHandler) AuthURL(c *gin.Context) {

}

// Callback 处理微信跳转回来的请求
func (w *WechatHandler) Callback(c *gin.Context) {

}
