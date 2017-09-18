package main

import (
	"github.com/openprovider/handlers/bearer"
	"github.com/openprovider/handlers/info"
	"github.com/takama/router"
)

func main() {
	r := router.New()

	version := "1.0.0"
	repo := "handlers"
	commit := "019cc819f8af4e2f7533fb3760f21387a4ef0cce"

	token := "mF_9.B5f-4.1JqM"
	h := func(ctx *router.Control) {}

	r.GET("/info", info.Handler(version, repo, commit))
	r.GET("/auth", bearer.WithToken(token, h))
	r.Listen(":3000")
}
