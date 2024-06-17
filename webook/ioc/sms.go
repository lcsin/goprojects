package ioc

import (
	"github.com/lcsin/goprojets/webook/internal/service/sms"
	"github.com/lcsin/goprojets/webook/internal/service/sms/local"
)

// InitSMSService 初始化短信服务
func InitSMSService() sms.Service {
	return local.NewService()
}
