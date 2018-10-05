package main

import (
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/icza/session"
	ghint "github.com/k8s-community/github-integration/client"
	umClient "github.com/k8s-community/user-manager/client"
	_ "github.com/lib/pq" // postgresql driver
	"github.com/openprovider/handlers/info"
	"github.com/takama/router"
	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/postgresql"

	"github.com/k8s-community/ui/handlers"
	"github.com/k8s-community/ui/session/storage"
	"github.com/k8s-community/ui/version"
)

var log logrus.Logger

func main() {

	// --use-service-discovery true|false

	var serviceDiscovery = flag.Bool("sd", true, "service discovery")
	flag.Parse()

	log := logrus.New()
	log.Formatter = new(logrus.TextFormatter)
	log.Level = logrus.DebugLevel
	logger := log.WithFields(logrus.Fields{"service": "ui"})

	var namespace string
	if *serviceDiscovery {
		var err error
		namespace, err = getFromEnv("NAMESPACE")
		if err != nil {
			logger.Fatalf("Namespace is not set: %+v", err)
		}
	}

	var errors []error

	dbUser, err := getFromEnv("UIDB_USER")
	if err != nil {
		errors = append(errors, err)
	}

	dbPass, err := getFromEnv("UIDB_PASSWORD")
	if err != nil {
		errors = append(errors, err)
	}

	dbName, err := getFromEnv("UIDB_NAME")
	if err != nil {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		logger.Fatalf("Couldn't start service because required DB parameters are not set: %+v", errors)
	}

	var dbHost, dbPort string

	if len(namespace) == 0 {
		var err error
		dbHost, err = getFromEnv("DB_HOST")

		if err != nil {
			logger.Fatalf("Couldn't start service %+v", err)
		}

		dbPort, err = getFromEnv("DB_PORT")

		if err != nil {
			logger.Fatalf("Couldn't start service %+v", err)
		}

	} else {
		dbHost = fmt.Sprintf("%s.%s", "uidb", namespace)
		dbPort = "5432"
	}

	db, err := startupDB(dbHost, dbPort, dbUser, dbPass, dbName)
	if err != nil {
		log.Fatalf("Couldn't start up DB for %v:%v: %+v", dbHost, dbPort, err)
	}

	// Session manager settings: temporary solution
	session.Global.Close()
	cookieMngrOptions := &session.CookieMngrOptions{
		SessIDCookieName: "k8s-community-session-id",
		AllowHTTP:        true,
		CookieMaxAge:     48 * time.Hour,
	}
	sessionStorage := storage.NewDB(db, logger)
	session.Global = session.NewCookieManagerOptions(sessionStorage, cookieMngrOptions)

	serviceHost, err := getFromEnv("SERVICE_HOST")
	if err != nil {
		errors = append(errors, err)
	}

	servicePort, err := getFromEnv("SERVICE_PORT")
	if err != nil {
		errors = append(errors, err)
	}

	usermanBaseURL := fmt.Sprintf("http://user-manager.%s:80", namespace)
	ghintBaseURL := fmt.Sprintf("http://github-integration.%s:80", namespace)

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

	// oauthState is a token to protect the user from CSRF attacks
	oauthState, err := getFromEnv("GITHUB_OAUTH_STATE")
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

	// Init github-integration client to get info about the builds
	ghintClient, err := ghint.NewClient(nil, ghintBaseURL)
	if err != nil {
		logger.Fatalf("Couldn't get an instance of github-integration's service client: %+v", err)
	}

	githubHandler := handlers.NewGitHubOAuth(logger, usermanClient, oauthState, githubClientID, githubClientSecret)

	// TODO: add graceful shutdown

	r := router.New()
	r.Handler("GET", "/static/*", http.StripPrefix("/static", http.FileServer(http.Dir("./static"))))
	r.GET("/", handlers.Home(db, logger, k8sGuestToken))
	r.GET("/oauth/github", githubHandler.Login)
	r.GET("/oauth/github-cb", githubHandler.Callback)
	r.GET("/signout", handlers.Signout())
	r.GET("/builds/:uuid", handlers.BuildHistory(ghintClient, "en"))

	r.GET("/builds/:id", func(c *router.Control) {
		c.Code(http.StatusOK).Body(http.StatusText(http.StatusOK))
	})

	r.GET("/info", info.Handler(version.RELEASE, version.REPO, version.COMMIT))
	r.GET("/healthz", func(c *router.Control) {
		c.Code(http.StatusOK).Body(http.StatusText(http.StatusOK))
	})

	r.NotFound = handlers.NotFound(logger)

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

// startupDB makes connection with DB, initializes reform DB level.
func startupDB(host, port, user, password, name string) (*reform.DB, error) {
	dataSource := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, name,
	)

	conn, err := sql.Open("postgres", dataSource)
	if err != nil {
		return nil, err
	}

	if err = conn.Ping(); err != nil {
		return nil, err
	}

	db := reform.NewDB(conn, postgresql.Dialect, reform.NewPrintfLogger(log.Printf))
	if err != nil {
		return nil, err
	}

	return db, nil
}

func startupDBWithConnectionString(connectionString string) (*reform.DB, error) {
	dataSource := fmt.Sprintf(connectionString)

	conn, err := sql.Open("postgres", dataSource)
	if err != nil {
		return nil, err
	}

	if err = conn.Ping(); err != nil {
		return nil, err
	}

	db := reform.NewDB(conn, postgresql.Dialect, reform.NewPrintfLogger(log.Printf))
	if err != nil {
		return nil, err
	}

	return db, nil
}
