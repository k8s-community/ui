package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// BuildCallback todo: add description
type BuildCallback struct {
	Username    string  `json:"username"`
	Repository  string  `json:"repository"`
	CommitHash  string  `json:"commitHash"`
	State       string  `json:"state"`
	BuildURL    *string `json:"buildURL,omitempty"`
	Description *string `json:"description,omitempty"`
	Context     *string `json:"context,omitempty"`
}

// CommitStatus todo: add description
type CommitStatus struct {
	State       string  `json:"state"`
	BuildURL    *string `json:"target_url,omitempty"`
	Description *string `json:"description,omitempty"`
	Context     *string `json:"context,omitempty"`
}

// UpdateCommitStatus todo: add description
func (c *Client) UpdateCommitStatus(build *BuildCallback) error {
	err := c.generateAccessToken()
	if err != nil {
		return fmt.Errorf("cannot generate access token: %s", err)
	}

	commitStatus := CommitStatus{
		State:       build.State,
		BuildURL:    build.BuildURL,
		Description: build.Description,
		Context:     build.Context,
	}

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(commitStatus)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/repos/%s/%s/statuses/%s", c.baseURL, build.Username, build.Repository, build.CommitHash), b)
	if err != nil {
		return fmt.Errorf("could not create request: %s", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("token %s", c.token.Token))

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("could not update commit status from GitHub API for %s//%s: %s", build.Username, build.Repository, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("received non 2xx response status %q when updating commit status: %s", resp.Status, req.URL)
	}

	if err := json.NewDecoder(resp.Body).Decode(&c.token); err != nil {
		return err
	}

	return nil
}
