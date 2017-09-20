package info_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/openprovider/handlers/info"
	"github.com/takama/router"
)

func TestInfo(t *testing.T) {
	c := new(router.Control)
	r := httptest.NewRequest("GET", "/", nil)
	c.Request = r

	w := httptest.NewRecorder()
	c.Writer = w

	version := "0.0.1"
	repo := "test"
	commit := "commit"

	h := info.Handler(version, repo, commit)
	h(c)

	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Response status is %d (expected %d)", resp.StatusCode, http.StatusOK)
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Content-Type s is %s (expected application/json)", contentType)
	}

	testResponse := struct {
		Version string `json:"version"`
		Repo    string `json:"repo"`
		Commit  string `json:"commit"`
	}{}

	err := json.NewDecoder(resp.Body).Decode(&testResponse)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if testResponse.Version != version {
		t.Errorf("Response version is %s (expected %s)", testResponse.Version, version)
	}

	if testResponse.Repo != repo {
		t.Errorf("Response repo is %s (expected %s)", testResponse.Repo, repo)
	}

	if testResponse.Commit != commit {
		t.Errorf("Response commit is %s (expected %s)", testResponse.Commit, commit)
	}
}
