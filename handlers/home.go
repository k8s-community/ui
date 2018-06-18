package handlers

import (
	"html/template"
	"net/http"

	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/icza/session"
	"github.com/k8s-community/ui/models"
	"github.com/takama/router"
	"gopkg.in/reform.v1"
)

// Home handles homepage request
func Home(db *reform.DB, log logrus.FieldLogger, k8sToken string) router.Handle {
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
			GitHubSignInLink string        // link to sign in to GitHub
			SignOutLink      string        // link to sign out (delete session)
			Login            string        // user's login
			Activated        bool          // is user activated in k8s
			GuestToken       string        // a token to reach Kubernetes
			Token            string        // personal token
			CA               template.HTML // personal cert
		}{
			GitHubSignInLink: "/oauth/github",
			SignOutLink:      "/signout",
			GuestToken:       k8sToken,
		}

		// Check if user have already logged in
		sessionData := session.Get(c.Request)
		if sessionData != nil {
			data.Login = sessionData.CAttr("Login").(string)
			data.Activated = sessionData.Attr("Activated").(bool)
		}

		token, cert := GetToken(db, log, data.Login)
		data.Token = token
		data.CA = template.HTML(cert)

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

func GetToken(db *reform.DB, logger logrus.FieldLogger, username string) (token string, cert string) {
	st, err := db.FindOneFrom(models.UserTable, "name", username)
	if err == reform.ErrNoRows {
		logger.Infof("Show user token and cert: attention! user '%s' not found", username)
		return
	}

	if err != nil {
		logger.Errorf("Show user token and cert: attention! Couldn't get user from DB: %+v", err)
		return
	}

	user := st.(*models.User)
	if user.Token != nil {
		token = *user.Token
	}
	if user.Cert != nil {
		cert = *user.Cert
		cert = strings.Replace(cert, "\n", "<br>", -1)
	}

	return token, cert
}
