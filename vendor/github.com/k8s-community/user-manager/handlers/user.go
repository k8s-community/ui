package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/k8s-community/user-manager/k8s"
	"github.com/takama/router"
)

// User defines
type User struct {
	Name string `json:"name"`
}

// SyncUser activates user in k8s system (creates namespaces, secrets)
func (h *Handler) SyncUser(c *router.Control) {
	var user User

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		h.Errlog.Printf("couldn't read request body: %s", err)
		c.Code(http.StatusBadRequest).Body(nil)
		return
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		h.Errlog.Printf("couldn't validate request body: %s", err)
		c.Code(http.StatusBadRequest).Body(nil)
		return
	}

	if len(user.Name) == 0 {
		c.Code(http.StatusBadRequest).Body(nil)
		return
	}

	h.Infolog.Printf("try to activate user %s", user.Name)

	client, err := k8s.NewClient("https://master.k8s.community:443", h.Env["K8S_TOKEN"])
	if err != nil {
		h.Errlog.Printf("cannot connect to k8s server: %s", err)
		c.Code(http.StatusInternalServerError).Body(nil)
		return
	}

	k8sUser := strings.ToLower(user.Name)
	namespace, _ := client.GetNamespace(k8sUser)
	if namespace != nil {
		h.Infolog.Printf("user %s already exists", k8sUser)
		c.Code(http.StatusOK).Body(nil)
		return
	}

	err = client.CreateNamespace(k8sUser)
	if err != nil {
		h.Errlog.Printf("%s", err)
		c.Code(http.StatusInternalServerError).Body(nil)
		return
	}

	secretNames := []string{h.Env["DOCKER_REGISTRY_SECRET_NAME"], h.Env["TLS_SECRET_NAME"]}
	for _, secretName := range secretNames {
		err = client.CopySecret(secretName, "default", k8sUser)
		if err != nil {
			h.Errlog.Printf("%s", err)
			c.Code(http.StatusInternalServerError).Body(nil)
			return
		}
	}

	err = client.CreateNamespaceAdmin(k8sUser)
	if err != nil {
		h.Errlog.Printf("%s", err)
		c.Code(http.StatusInternalServerError).Body(nil)
		return
	}

	tok, err := client.GetNamespaceToken(k8sUser)
	if err != nil {
		h.Errlog.Printf("%s", err)
		c.Code(http.StatusInternalServerError).Body(nil)
		return
	}

	h.Infolog.Printf("the token is: %v", tok)
	c.Code(http.StatusOK).Body(tok)

	h.Infolog.Printf("user %s is activated", k8sUser)
}
