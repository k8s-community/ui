package handlers

import (
	"log"
	"net/http"

	"github.com/k8s-community/user-manager/version"
	"github.com/takama/router"
)

// Handler defines
type Handler struct {
	Infolog *log.Logger
	Errlog  *log.Logger
	Env     map[string]string
}

func (h *Handler) HealthzHandler(c *router.Control) {
	c.Code(http.StatusOK).Body("Ok")
}

func (h *Handler) InfoHandler(c *router.Control) {
	c.Code(http.StatusOK).Body(
		map[string]string{
			"version": version.RELEASE,
			"commit":  version.COMMIT,
			"repo":    version.REPO,
		},
	)
}
