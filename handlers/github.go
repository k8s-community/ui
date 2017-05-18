package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	client "github.com/google/go-github/github"
	"github.com/icza/session"
	"github.com/takama/router"
	"golang.org/x/oauth2"
	oAuth "golang.org/x/oauth2/github"
)

// GitHubOAuth is a handler set to use GitHubOAuth features
type GitHubOAuth struct {
	state     string
	oAuthConf *oauth2.Config
	log       logrus.FieldLogger
}

// NewGitHubOAuth create new GitHubOAuth handler set:
// - state is a token to protect the user from CSRF attacks
// - clientID and clientSecret are the parameters from github.com/settings/developers
func NewGitHubOAuth(log logrus.FieldLogger, state, clientID, clientSecret string) *GitHubOAuth {
	conf := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{},
		Endpoint:     oAuth.Endpoint,
	}

	return &GitHubOAuth{
		state:     state,
		oAuthConf: conf,
		log:       log,
	}
}

// Login is a handler to redirect to GitHub authorization page
func (h *GitHubOAuth) Login(c *router.Control) {
	url := h.oAuthConf.AuthCodeURL(h.state, oauth2.AccessTypeOnline)
	http.Redirect(c.Writer, c.Request, url, http.StatusTemporaryRedirect)
}

// Callback is a handler to process authorization callback from GitHub
func (h *GitHubOAuth) Callback(c *router.Control) {
	state := c.Get("state")
	code := c.Get("code")

	if state != h.state {
		h.log.Errorf("Wrong state %s with code %s", state, code)
		http.Redirect(c.Writer, c.Request, "/", http.StatusMovedPermanently)
		return
	}

	ctx := context.Background()
	token, err := h.oAuthConf.Exchange(ctx, code)

	if err != nil {
		h.log.Errorf("Exchange failed for code %s: %+v", code, err)
		http.Redirect(c.Writer, c.Request, "/", http.StatusMovedPermanently)
		return
	}

	oauthClient := h.oAuthConf.Client(ctx, token)
	githubClient := client.NewClient(oauthClient)
	user, _, err := githubClient.Users.Get(ctx, "")
	if err != nil || user.Login == nil {
		h.log.Errorf("Couldn't get user for code %s: %+v", code, err)
		http.Redirect(c.Writer, c.Request, "/", http.StatusMovedPermanently)
		return
	}

	h.log.WithField("user", *user.Login).Info("GitHub user was authorized in oauth-proxy")

	sessionData := session.NewSessionOptions(&session.SessOptions{
		CAttrs: map[string]interface{}{"Login": *user.Login},
		Attrs:  map[string]interface{}{"Activated": false},
	})
	session.Add(sessionData, c.Writer)
	userFields := logrus.Fields{"user": *user.Login, "session": sessionData.ID()}
	h.log.WithFields(userFields).Infof("Session was created")

	go func() {
		// Call User-Manager instead of sleep
		time.Sleep(5 * time.Second)

		h.log.WithField("user", *user.Login).Info("GitHub user was activated in Kubernetes")
		sessionData.SetAttr("Activated", true)
		h.log.WithFields(userFields).Infof("Session was updated: set 'activated' value")
	}()

	http.Redirect(c.Writer, c.Request, "/", http.StatusMovedPermanently)
}
