package handlers

import (
	"html/template"

	"github.com/Sirupsen/logrus"
	"github.com/icza/session"
	"github.com/takama/router"
)

// Home handles homepage request
func Home(log logrus.FieldLogger, k8sToken string) router.Handle {
	return func(c *router.Control) {
		t, err := template.ParseFiles("templates/layout.html", "templates/index.html")

		if err != nil {
			log.Fatalf("Couldn't parse template files: %+v", err)
		}

		data := struct {
			GitHubSignInLink string // link to sign in to GitHub
			Login            string // user's login
			Activated        bool   // is user activated in k8s
			GuestToken       string // a token to reach Kubernetes
		}{
			GitHubSignInLink: "/oauth/github",
			GuestToken:       k8sToken,
		}

		// Check if user have already logged in
		sessionData := session.Get(c.Request)
		if sessionData != nil {
			data.Login = sessionData.CAttr("Login").(string)
			data.Activated = sessionData.Attr("Activated").(bool)
		}

		t.ExecuteTemplate(c.Writer, "layout", data)
	}
}
