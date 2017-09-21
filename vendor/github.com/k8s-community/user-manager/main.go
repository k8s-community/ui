package main

import (
	"log"
	"os"

	"os/signal"
	"syscall"

	"github.com/k8s-community/user-manager/handlers"
	"github.com/takama/router"
)

const (
	apiPrefix = "/api/v1"
)

// main function
func main() {
	keys := []string{
		"USERMAN_SERVICE_PORT",
		"DOCKER_REGISTRY_SECRET_NAME", "TLS_SECRET_NAME",
		"K8S_HOST", "K8S_TOKEN",
	}
	h := &handlers.Handler{
		Infolog: log.New(os.Stdout, "[USERMAN:INFO]: ", log.LstdFlags),
		Errlog:  log.New(os.Stderr, "[USERMAN:ERROR]: ", log.LstdFlags),
		Env:     make(map[string]string, len(keys)),
	}
	for _, key := range keys {
		value := os.Getenv(key)
		if value == "" {
			h.Errlog.Fatalf("%s environment variable was not set", key)
		}
		h.Env[key] = value
	}

	r := router.New()

	r.GET("/healthz", h.HealthzHandler)
	r.GET("/info", h.InfoHandler)

	r.PUT(apiPrefix+"/sync-user", h.SyncUser)

	h.Infolog.Printf("start listening port %s", h.Env["USERMAN_SERVICE_PORT"])
	h.Infolog.Printf("registered routes are: %+v", r.Routes())

	go r.Listen(":" + h.Env["USERMAN_SERVICE_PORT"])

	// Set up channel on which to send signal notifications.
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal is sent.
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)
	killSignal := <-interrupt
	h.Infolog.Println("Got signal:", killSignal)

	if killSignal == os.Kill {
		h.Infolog.Println("Service was killed")
	} else {
		h.Infolog.Println("Service was terminated by system signal")
	}

	h.Infolog.Println("shutdown")
}
