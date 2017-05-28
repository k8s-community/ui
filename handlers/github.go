package handlers

import (
	"context"
	"net/http"

	"github.com/Sirupsen/logrus"
	ghClient "github.com/google/go-github/github"
	"github.com/icza/session"
	umClient "github.com/k8s-community/user-manager/client"
	"github.com/takama/router"
	"golang.org/x/oauth2"
	ghOAuth "golang.org/x/oauth2/github"
)

// GitHubOAuth is a handler set to use GitHubOAuth features
type GitHubOAuth struct {
	state         string
	oAuthConf     *oauth2.Config
	log           logrus.FieldLogger
	usermanClient *umClient.Client
}

// NewGitHubOAuth create new GitHubOAuth handler set:
// - state is a token to protect the user from CSRF attacks
// - clientID and clientSecret are the parameters from github.com/settings/developers
func NewGitHubOAuth(log logrus.FieldLogger, umClient *umClient.Client, state, ghClientID, ghClientSecret string) *GitHubOAuth {
	conf := &oauth2.Config{
		ClientID:     ghClientID,
		ClientSecret: ghClientSecret,
		Endpoint:     ghOAuth.Endpoint,
	}

	return &GitHubOAuth{
		state:         state,
		oAuthConf:     conf,
		log:           log,
		usermanClient: umClient,
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
	githubClient := ghClient.NewClient(oauthClient)
	user, _, err := githubClient.Users.Get(ctx, "")
	if err != nil || user.Login == nil {
		h.log.Errorf("Couldn't get user for code %s: %+v", code, err)
		http.Redirect(c.Writer, c.Request, "/", http.StatusMovedPermanently)
		return
	}

	h.log.WithField("user", *user.Login).Info("GitHub user was authorized in oauth-proxy")

	sessionData := session.NewSessionOptions(&session.SessOptions{
		CAttrs: map[string]interface{}{"Login": *user.Login},
		Attrs:  map[string]interface{}{"Activated": false, "HasError": false},
	})
	session.Add(sessionData, c.Writer)

	go h.syncUser(*user.Login, sessionData)

	http.Redirect(c.Writer, c.Request, "/", http.StatusMovedPermanently)
}

func (h *GitHubOAuth) syncUser(login string, sessionData session.Session) {
	logger := h.log.WithFields(logrus.Fields{"user": login, "session": sessionData.ID()})
	logger.Infof("Session was created")

	user := umClient.NewUser(login)
	status, err := h.usermanClient.User.Sync(user)

	if err != nil {
		logger.Info("Error during user Kubernetes sync: %+v", err)
		sessionData.SetAttr("Activated", false)
		sessionData.SetAttr("HasError", true)
		return
	}

	logger.Infof("Status from user-manager service is: %d", status)

	sessionData.SetAttr("Activated", true)
	sessionData.SetAttr("HasError", false)
	logger.Infof("Session was updated: set 'activated' value")
}
