package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/k8s-community/github-integration/handlers"
	_ "github.com/lib/pq" // postgresql driver
	"github.com/takama/router"
	"gopkg.in/reform.v1"
	"gopkg.in/reform.v1/dialects/postgresql"
)

const (
	apiPrefix = "/api/v1"
)

// main function
func main() {
	var errors []error

	// Database settings
	namespace, err := getFromEnv("POD_NAMESPACE")
	if err != nil {
		errors = append(errors, err)
	}

	dbUser, err := getFromEnv("GITHUBDB_USER")
	if err != nil {
		errors = append(errors, err)
	}

	dbPass, err := getFromEnv("GITHUBDB_PASSWORD")
	if err != nil {
		errors = append(errors, err)
	}

	dbName, err := getFromEnv("GITHUBDB_NAME")
	if err != nil {
		errors = append(errors, err)
	}

	dbHost := fmt.Sprintf("%s.%s", "db-github", namespace)
	dbPort := "5432"

	db, err := startupDB(dbHost, dbPort, dbUser, dbPass, dbName)
	if err != nil {
		log.Fatalf("Couldn't start up DB: %+v", err)
	}

	if len(errors) > 0 {
		log.Fatalf("Couldn't start service because required DB parameters are not set: %+v", errors)
	}

	keys := []string{
		"GITHUBINT_LOCAL_PORT", "GITHUBINT_BRANCH",
		"GITHUBINT_TOKEN", "GITHUBINT_PRIV_KEY", "GITHUBINT_INTEGRATION_ID",
		"USERMAN_BASE_URL", "CICD_BASE_URL",
	}

	h := &handlers.Handler{
		DB:      db,
		Infolog: log.New(os.Stdout, "[GITHUBINT:INFO]: ", log.LstdFlags),
		Errlog:  log.New(os.Stderr, "[GITHUBINT:ERROR]: ", log.LstdFlags),
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
	r.PanicHandler = handlers.Panic

	r.GET(apiPrefix+"/", h.HomeHandler)

	r.GET("/healthz", h.HealthzHandler)
	r.GET("/info", h.InfoHandler)

	r.GET(apiPrefix+"/home", h.HomeHandler)
	r.POST(apiPrefix+"/webhook", h.WebHookHandler)
	r.POST(apiPrefix+"/auth-callback", h.AuthCallbackHandler)
	r.POST(apiPrefix+"/build-cb", h.BuildCallbackHandler)
	r.POST(apiPrefix+"/build-results", h.BuildResultsHandler)
	r.GET(apiPrefix+"/build-results/:uuid", h.ShowBuildResults)
	h.Infolog.Printf("start listening port %s", h.Env["GITHUBINT_LOCAL_PORT"])
	h.Infolog.Printf("Registered routes are: %+v", r.Routes())

	go r.Listen(":" + h.Env["GITHUBINT_LOCAL_PORT"])

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
