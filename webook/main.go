package main

import (
	"fmt"

	"github.com/lcsin/gopocket/util/httpx"
	"github.com/lcsin/goprojets/webook/config"
)

func main() {
	r := InitWebServer()
	httpx.Graceful(r, fmt.Sprintf(":%v", config.Cfg.Port))
}
