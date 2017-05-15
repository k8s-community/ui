package main

import (
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/k8s-community/oauth-proxy/handlers"
	"github.com/satori/go.uuid"
	"github.com/takama/router"
)

func main() {
	log := logrus.New()

	var errors []error

	serviceHost, err := getFromEnv("SERVICE_HOST")
	if err != nil {
		errors = append(errors, err)
	}

	servicePort, err := getFromEnv("SERVICE_PORT")
	if err != nil {
		errors = append(errors, err)
	}

	githubClientID, err := getFromEnv("GITHUB_CLIENT_ID")
	if err != nil {
		errors = append(errors, err)
	}

	githubClientSecret, err := getFromEnv("GITHUB_CLIENT_SECRET")
	if err != nil {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		log.Fatalf("Couldn't start service because required parameters are not set: %+v", errors)
	}

	// oauthState is a token to protect the user from CSRF attacks
	oauthState := uuid.NewV4().String()

	githubHandler := handlers.NewGitHubOAuth(oauthState, githubClientID, githubClientSecret)

	// TODO: add graceful shutdown

	r := router.New()
	r.GET("/oauth/login", githubHandler.Login)
	r.GET("/oauth/github-cb", githubHandler.Callback)

	hostPort := fmt.Sprintf("%s:%s", serviceHost, servicePort)
	log.Infof("Ready to listen %s\nRoutes: %+v", hostPort, r.Routes())
	r.Listen(hostPort)
}

func getFromEnv(name string) (string, error) {
	value := os.Getenv(name)
	if len(value) == 0 {
		return "", fmt.Errorf("Environement variable %s must be set", name)
	}

	return value, nil
}
