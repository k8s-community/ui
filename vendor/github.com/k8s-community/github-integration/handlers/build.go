package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/k8s-community/github-integration/github"
	"github.com/takama/router"
)

// BuildCallbackHandler is handler for callback from build service (system)
func (h *Handler) BuildCallbackHandler(c *router.Control) {
	h.Infolog.Print("Received callback request...")

	var build github.BuildCallback

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		h.Errlog.Printf("couldn't read request body: %s", err)
		c.Code(http.StatusBadRequest).Body(nil)
		return
	}

	err = json.Unmarshal(body, &build)
	if err != nil {
		h.Errlog.Printf("couldn't validate request body: %s", err)
		c.Code(http.StatusBadRequest).Body(nil)
		return
	}

	err = h.updateCommitStatus(c, &build)
	if err != nil {
		h.Errlog.Printf("cannot update commit status, build: %+v, err: %s", build, err)
		return
	}

	c.Code(http.StatusOK).Body(nil)
}
