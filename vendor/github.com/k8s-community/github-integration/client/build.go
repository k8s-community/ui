package client

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	buildCallbackURLStr = "/build-cb"
	buildCResultsURLStr = "/build-results"
)

// Possible GitHub Build states
const (
	StatePending = "pending"
	StateSuccess = "success"
	StateError   = "error"
	StateFailure = "failure"
)

const (
	ContextCICD = "k8s-community/cicd"
)

// BuildService defines
type BuildService struct {
	client *Client
}

// BuildCallback defines
type BuildCallback struct {
	Username    string `json:"username"`
	Repository  string `json:"repository"`
	CommitHash  string `json:"commitHash"`
	State       string `json:"state"`
	BuildURL    string `json:"buildURL"`
	Description string `json:"description"`
	Context     string `json:"context"`
}

type BuildResults struct {
	UUID       string `json:"uuid"`
	Username   string `json:"username"`
	Repository string `json:"repository"`
	CommitHash string `json:"commitHash"`
	Passed     bool   `json:"passed"`
	Log        string `json:"log"`
}

// BuildCallback sends request for update commit status on github side
func (u *BuildService) BuildCallback(build BuildCallback) error {
	req, err := u.client.NewRequest(postMethod, buildCallbackURLStr, build)
	if err != nil {
		return err
	}

	_, err = u.client.Do(req, nil)
	if err != nil {
		requestBody, _ := json.Marshal(build)
		return fmt.Errorf("couldn't process build-callback request: %v, request body:'%s'", err, requestBody)
	}

	return nil
}

func (u *BuildService) BuildResults(results *BuildResults) error {
	req, err := u.client.NewRequest(postMethod, buildCResultsURLStr, results)
	if err != nil {
		return err
	}

	_, err = u.client.Do(req, nil)
	if err != nil {
		requestBody, _ := json.Marshal(results)
		return fmt.Errorf("couldn't process build-results request: %v, request body:'%s'", err, requestBody)
	}

	return nil
}

func (u *BuildService) ShowResults(uuid string) (*BuildResults, error) {
	req, err := u.client.NewRequest(http.MethodGet, buildCResultsURLStr+"/"+uuid, nil)
	if err != nil {
		return nil, err
	}

	build := &BuildResults{}
	_, err = u.client.Do(req, build)
	if err != nil {
		return nil, fmt.Errorf("couldn't get results for uuid %s: %v", uuid, err)
	}

	return build, nil
}
