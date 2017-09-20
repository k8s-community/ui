package bearer_test

import (
	"github.com/openprovider/handlers/bearer"
	"github.com/takama/router"
)

func ExampleWithToken() {
	r := router.New()

	token := "s-fdF8-mF_9.B-4.1Cfd"
	h := func(ctx *router.Control) {}

	r.GET("/test", bearer.WithToken(token, h))
	r.Listen(":3000")
}
