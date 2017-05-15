package handlers

import (
	"context"
	"fmt"
	"net/http"

	client "github.com/google/go-github/github"
	"github.com/takama/router"
	"golang.org/x/oauth2"
	oAuth "golang.org/x/oauth2/github"
)

// GitHubOAuth is a handler set to use GitHubOAuth features
type GitHubOAuth struct {
	State     string
	OAuthConf *oauth2.Config
}

// NewGitHubOAuth create new GitHubOAuth handler set:
// - state is a token to protect the user from CSRF attacks
// - clientID and clientSecret are the parameters from github.com/settings/developers
func NewGitHubOAuth(state, clientID, clientSecret string) *GitHubOAuth {
	conf := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{"user"},
		Endpoint:     oAuth.Endpoint,
	}

	return &GitHubOAuth{
		State:     state,
		OAuthConf: conf,
	}
}

// Login is a handler to redirect to GitHub authorization page
func (h *GitHubOAuth) Login(c *router.Control) {
	url := h.OAuthConf.AuthCodeURL(h.State, oauth2.AccessTypeOnline)
	http.Redirect(c.Writer, c.Request, url, http.StatusTemporaryRedirect)
}

// Callback is a handler to process authorization callback from GitHub
func (h *GitHubOAuth) Callback(c *router.Control) {
	state := c.Get("state")
	if state != state {
		// TODO: log invalid state
		return
	}

	code := c.Get("code")
	ctx := context.Background()

	token, err := h.OAuthConf.Exchange(ctx, code)
	if err != nil {
		// TODO: Log exchange failed
		return
	}

	oauthClient := h.OAuthConf.Client(ctx, token)
	githubClient := client.NewClient(oauthClient)
	user, _, err := githubClient.Users.Get(ctx, "")
	if err != nil {
		// TODO: Log 'can't get user'
		return
	}

	// TODO: log what user was logged in

	go func() {
		// Call User-Manager here
		fmt.Println(user)
	}()
}
