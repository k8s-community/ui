package info_test

import (
	"github.com/openprovider/handlers/info"
	"github.com/takama/router"
)

// ExampleHandler is a usage example for info.Handler
func ExampleHandler() {
	r := router.New()

	version := "1.0.0"
	repo := "handlers"
	commit := "019cc819f8af4e2f7533fb3760f21387a4ef0cce"

	r.GET("/info", info.Handler(version, repo, commit))
	r.Listen(":3000")
}
