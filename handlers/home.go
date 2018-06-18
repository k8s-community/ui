package handlers

import (
	"html/template"

	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/icza/session"
	"github.com/takama/router"
)

// Home handles homepage request
func Home(log logrus.FieldLogger, k8sToken string) router.Handle {
	lang := "en"
	return func(c *router.Control) {
		t, err := template.ParseFiles(
			"templates/"+lang+"/layout.html",
			"templates/"+lang+"/index.html",
		)

		if err != nil {
			log.Fatalf("Couldn't parse template files: %+v", err)
		}

		data := struct {
			GitHubSignInLink string // link to sign in to GitHub
			SignOutLink      string // link to sign out (delete session)
			Login            string // user's login
			Activated        bool   // is user activated in k8s
			GuestToken       string // a token to reach Kubernetes
			Token            string // personal token
			CA               string // personal cert
		}{
			GitHubSignInLink: "/oauth/github",
			SignOutLink:      "/signout",
			GuestToken:       k8sToken,
		}

		// Check if user have already logged in
		sessionData := session.Get(c.Request)
		if sessionData != nil {
			data.Login = sessionData.CAttr("Login").(string)
			//data.Token = sessionData.CAttr("Token").(string)
			//data.CA = sessionData.CAttr("CA").(string)
			data.Activated = sessionData.Attr("Activated").(bool)
		}

		t.ExecuteTemplate(c.Writer, "layout", data)
	}
}

func Signout() router.Handle {
	return func(c *router.Control) {
		session.Remove(session.Get(c.Request), c.Writer)
		http.Redirect(c.Writer, c.Request, "/", http.StatusFound)
	}
}

// Handle undefined routes
func NotFound(log logrus.FieldLogger) router.Handle {
	return func(c *router.Control) {
		log.Warningf("couldn't find path: %s", c.Request.RequestURI)
		http.NotFound(c.Writer, c.Request)
	}
}
