package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/icza/session"
	common_handlers "github.com/k8s-community/handlers"
	"github.com/k8s-community/k8s-community/handlers"
	"github.com/k8s-community/k8s-community/version"
	umClient "github.com/k8s-community/user-manager/client"
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

	usermanBaseURL, err := getFromEnv("USERMAN_BASE_URL")
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

	k8sGuestToken, err := getFromEnv("K8S_GUEST_TOKEN")
	if err != nil {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		logger.Fatalf("Couldn't start service because required parameters are not set: %+v", errors)
	}

	// Init user-manager client to be able to create user in Kubernetes
	usermanClient, err := umClient.NewClient(nil, usermanBaseURL)
	if err != nil {
		logger.Fatalf("Couldn't get an instance of user-manager's service client: %+v", err)
	}

	// oauthState is a token to protect the user from CSRF attacks
	oauthState := uuid.NewV4().String()

	githubHandler := handlers.NewGitHubOAuth(logger, usermanClient, oauthState, githubClientID, githubClientSecret)

	// TODO: add graceful shutdown

	r := router.New()
	r.Handler("GET", "/static/*", http.StripPrefix("/static", http.FileServer(http.Dir("./static"))))
	r.GET("/", handlers.Home(logger, k8sGuestToken))
	r.GET("/oauth/github", githubHandler.Login)
	r.GET("/oauth/github-cb", githubHandler.Callback)

	r.GET("/info", func(c *router.Control) {
		common_handlers.Info(c, version.RELEASE, version.REPO, version.COMMIT)
	})
	r.GET("/healthz", func(c *router.Control) {
		c.Code(http.StatusOK).Body(http.StatusText(http.StatusOK))
	})

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
