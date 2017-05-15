package handlers

import (
	"html/template"

	"github.com/Sirupsen/logrus"
	"github.com/takama/router"
)

// Home handles homepage request
func Home(log logrus.FieldLogger) router.Handle {
	return func(c *router.Control) {
		t, err := template.ParseFiles("templates/layout.html", "templates/index.html")

		if err != nil {
			log.Fatal("Couldn't parse template files")
		}

		data := struct {
			GitHubSignInLink string
		}{
			GitHubSignInLink: "/oauth/github",
		}

		t.ExecuteTemplate(c.Writer, "layout", data)
	}
}
