package main

import (
	"github.com/lcsin/gopocket/util/httpx"
	"github.com/lcsin/goprojets/webook/internal/web"
)

func main() {
	r := web.RegisterRoutes()
	httpx.Graceful(r, ":8080")
}
