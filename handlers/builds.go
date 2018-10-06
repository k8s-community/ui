package handlers

import (
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/icza/session"
	ghint "github.com/k8s-community/github-integration/client"
	"github.com/takama/router"
)

func BuildHistory(client *ghint.Client, lang string) router.Handle {
	t, err := template.ParseFiles(
		"templates/"+lang+"/layout.html",
		"templates/"+lang+"/build-results.html",
	)
	if err != nil {
		log.Fatalf("Couldn't parse template files: %+v", err)
	}

	return func(c *router.Control) {
		sessionData := session.Get(c.Request)
		if sessionData == nil {
			http.Redirect(c.Writer, c.Request, "/", http.StatusFound)
			return
		}

		uuid := c.Get(":uuid")
		build, err := client.Build.ShowResults(uuid)
		if err != nil {
			log.Printf("couldn't get build logs: %v", err)
		}

		build.Log = strings.Replace(build.Log, "\n", "<br>", -1)

		t.ExecuteTemplate(c.Writer, "layout", template.HTML(build.Log))
	}
}
