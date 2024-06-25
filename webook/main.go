package main

import (
	"fmt"

	"github.com/lcsin/gopocket/util/httpx"
	"github.com/lcsin/goprojets/webook/ioc"
	"github.com/spf13/viper"
)

func main() {
	ioc.InitLocalConfig()
	r := InitWebServer()
	httpx.Graceful(r, fmt.Sprintf(":%v", viper.Get("service.port")))
}
