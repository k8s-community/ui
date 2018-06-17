package handlers

import (
	"log"
	"runtime/debug"

	"github.com/takama/router"
)

func Panic(c *router.Control) {
	log.Printf("Recovered panic:\n%s\n", string(debug.Stack()))
}
