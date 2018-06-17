package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// acceptHeader is the GitHub Integrations Preview Accept header.
const (
	acceptHeader = "application/vnd.github.machine-man-preview+json"
	apiBaseURL   = "https://api.github.com"
)

// Client definess
type Client struct {
	// HTTP client used to communicate with the API.
	client *http.Client

	// Base URL for API requests.
	baseURL *url.URL

	integrationID  int
	installationID int

	privKey []byte

	token *accessToken // token is the installation's access token
}

// accessToken is an installation access token response from GitHub
type accessToken struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

// NewClient initializes a Client instance
func NewClient(httpClient *http.Client, integrationID int, installationID int, privKey []byte) (*Client, error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	baseURL, err := url.Parse(apiBaseURL)
	if err != nil {
		return nil, fmt.Errorf("github client: cannot parse url %s: %s", apiBaseURL, err)
	}

	c := &Client{
		client:         httpClient,
		baseURL:        baseURL,
		installationID: installationID,
		integrationID:  integrationID,
		privKey:        privKey,
	}

	return c, nil
}

// NewRequest creates new http.Request instance
func (c *Client) NewRequest(method string, urlStr string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("cannot parse url %s: %s", urlStr, err)
	}

	u := c.baseURL.ResolveReference(rel)

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, fmt.Errorf("cannot encode data: %s", err)
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, fmt.Errorf("cannot send request: %s", err)
	}

	return req, nil
}

// Response is an API response.
// This wraps the standard http.Response.
type Response struct {
	*http.Response
}

// newResponse creates a new Response for the provided http.Response.
func newResponse(r *http.Response) *Response {
	response := &Response{Response: r}
	return response
}

// Do sends an API request and returns the API response.  The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error if an API error has occurred.  If v implements the io.Writer
// interface, the raw response body will be written to v, without attempting to
// first decode it.
func (c *Client) Do(req *http.Request, v interface{}) (*Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	response := newResponse(resp)

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
			if err == io.EOF {
				err = nil // ignore EOF errors caused by empty response body
			}
		}
	}

	return response, err
}

// generateBearer is used for JWT token generation
func (c *Client) generateBearer() (string, error) {
	parsedKey, err := jwt.ParseRSAPrivateKeyFromPEM(c.privKey)
	if err != nil {
		return "", err
	}

	bearer := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iat": time.Now().Unix(),
		"exp": time.Now().UTC().Add(time.Minute * time.Duration(10)).Unix(),
		"iss": c.integrationID,
	})

	// Sign and get the complete encoded token as a string using the secret
	bearerString, err := bearer.SignedString(parsedKey)

	return bearerString, err
}

// generateAccessToken is used for access token generation
func (c *Client) generateAccessToken() error {
	bearer, err := c.generateBearer()
	if err != nil {
		return fmt.Errorf("cannot generate bearer token: %s", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/installations/%v/access_tokens", c.baseURL, c.installationID), nil)
	if err != nil {
		return fmt.Errorf("could not create request: %s", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", bearer))
	req.Header.Set("Accept", acceptHeader)

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("could not get access_tokens from GitHub API for installation ID %v: %s", c.installationID, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("received non 2xx response status %q when fetching %v", resp.Status, req.URL)
	}

	if err := json.NewDecoder(resp.Body).Decode(&c.token); err != nil {
		return err
	}

	return nil
}
