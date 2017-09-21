package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

const (
	postMethod   = "POST"
	getMethod    = "GET"
	putMethod    = "PUT"
	deleteMethod = "DELETE"
)

const (
	apiPrefix = "/api/v1"
)

// Client defines
type Client struct {
	// HTTP client used to communicate with the API.
	client *http.Client

	// Base URL for API requests.
	BaseURL *url.URL

	// Services used for talking to different parts of the API.
	User *UserService
}

// NewClient creates a new Client instance
func NewClient(httpClient *http.Client, strBaseURL string) (*Client, error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	baseURL, err := url.Parse(strBaseURL)
	if err != nil {
		return nil, fmt.Errorf("user manager client: cannot parse url %s: %s", strBaseURL, err)
	}

	c := &Client{
		client:  httpClient,
		BaseURL: baseURL,
	}

	c.User = &UserService{client: c}

	return c, nil
}

// NewRequest creates a new http.Request instance
func (c *Client) NewRequest(method string, urlStr string, body interface{}) (*http.Request, error) {
	u, err := url.Parse(c.BaseURL.String() + apiPrefix + urlStr)
	if err != nil {
		return nil, fmt.Errorf("cannot parse url %s: %s", apiPrefix+urlStr, err)
	}

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, fmt.Errorf("cannot encode data: %s", err)
		}
	}

	// TODO: add better logger here
	log.Printf("Send %s request to %s\n", method, u.String())
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
