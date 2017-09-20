package bearer_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/openprovider/handlers/bearer"
	"github.com/takama/router"
)

func TestWithToken(t *testing.T) {
	token := "mF_9.B5f-4.1JqM"

	for _, c := range positiveSet(token, t) {
		w := httptest.NewRecorder()
		c.Writer = w

		h := func(ctx *router.Control) { ctx.Writer.WriteHeader(http.StatusOK) }
		h = bearer.WithToken(token, h)
		h(&c)

		resp := w.Result()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Response status is %d (expected %d)", resp.StatusCode, http.StatusOK)
		}
	}

	for _, c := range negativeSet(token, t) {
		w := httptest.NewRecorder()
		c.Writer = w

		h := func(ctx *router.Control) { ctx.Writer.WriteHeader(http.StatusOK) }
		h = bearer.WithToken(token, h)
		h(&c)

		resp := w.Result()
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Response status is %d (expected %d)", resp.StatusCode, http.StatusBadRequest)
		}
	}
}

// set of positive cases for test
func positiveSet(token string, t *testing.T) []router.Control {
	var testSet []router.Control

	// from header
	c := new(router.Control)
	r, err := http.NewRequest("GET", "test", nil)
	if err != nil {
		t.Error(err)
	}
	r.Header.Set("Authorization", "Bearer "+token)
	c.Request = r
	testSet = append(testSet, *c)

	// from body
	c = new(router.Control)
	data := url.Values{}
	data.Set("access_token", token)
	r, err = http.NewRequest("POST", "test", bytes.NewBufferString(data.Encode()))
	if err != nil {
		t.Error(err)
	}
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Request = r
	testSet = append(testSet, *c)

	// from query string
	r, err = http.NewRequest("GET", "test", nil)
	if err != nil {
		t.Error(err)
	}
	c.Request = r

	param := router.Param{
		Key:   "access_token",
		Value: token,
	}
	c.Set(param)
	testSet = append(testSet, *c)

	return testSet
}

func negativeSet(token string, t *testing.T) []router.Control {
	var testSet []router.Control

	// from body if content type is wrong
	c := new(router.Control)
	data := url.Values{}
	data.Set("access_token", token)
	r, err := http.NewRequest("POST", "test", bytes.NewBufferString(data.Encode()))
	if err != nil {
		t.Error(err)
	}
	r.Header.Set("Content-Type", "application/json")
	c.Request = r
	testSet = append(testSet, *c)

	// if token is not set at all
	c = new(router.Control)
	r, err = http.NewRequest("GET", "test", nil)
	if err != nil {
		t.Error(err)
	}
	c.Request = r
	testSet = append(testSet, *c)

	// wrong token from header
	c = new(router.Control)
	r, err = http.NewRequest("GET", "test", nil)
	if err != nil {
		t.Error(err)
	}
	r.Header.Set("Authorization", "Bearer abc"+token)
	c.Request = r
	testSet = append(testSet, *c)

	return testSet
}
