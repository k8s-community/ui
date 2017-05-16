package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/icza/session"
	"github.com/k8s-community/oauth-proxy/handlers"
	"github.com/satori/go.uuid"
	"github.com/takama/router"
)

func main() {
	log := logrus.New()
	log.Formatter = new(logrus.TextFormatter)
	logger := log.WithFields(logrus.Fields{"service": "oauth-proxy"})

	// Session manager settings: temporary solution
	session.Global.Close()
	cookieMngrOptions := &session.CookieMngrOptions{
		SessIDCookieName: "k8s-community-session-id",
		AllowHTTP:        true,
		CookieMaxAge:     48 * time.Hour,
	}
	session.Global = session.NewCookieManagerOptions(session.NewInMemStore(), cookieMngrOptions)

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
		logger.Fatalf("Couldn't start service because required parameters are not set: %+v", errors)
	}

	// oauthState is a token to protect the user from CSRF attacks
	oauthState := uuid.NewV4().String()

	githubHandler := handlers.NewGitHubOAuth(logger, oauthState, githubClientID, githubClientSecret)

	// TODO: add graceful shutdown

	r := router.New()
	r.Handler("GET", "/static/*", http.StripPrefix("/static", http.FileServer(http.Dir("./static"))))
	r.GET("/", handlers.Home(logger))
	r.GET("/oauth/github", githubHandler.Login)
	r.GET("/oauth/github-cb", githubHandler.Callback)

	hostPort := fmt.Sprintf("%s:%s", serviceHost, servicePort)
	logger.Infof("Ready to listen %s\nRoutes: %+v", hostPort, r.Routes())
	r.Listen(hostPort)
}

func getFromEnv(name string) (string, error) {
	value := os.Getenv(name)
	if len(value) == 0 {
		return "", fmt.Errorf("Environement variable %s must be set", name)
	}

	return value, nil
}
