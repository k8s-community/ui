package bearer

import (
	"bytes"
	"net/http"
	"net/url"
	"testing"

	"github.com/takama/router"
)

func TestFromHeader(t *testing.T) {
	token := "mF_9.B5f-4.1JqM"

	// positive case
	r, err := http.NewRequest("GET", "test", nil)
	if err != nil {
		t.Error(err)
	}
	r.Header.Set("Authorization", "Bearer "+token)

	res, err := fromHeader(r)
	if err != nil {
		t.Error(err)
	}
	if res != token {
		t.Errorf("Token %s is wrong, expected: %s", res, token)
	}

	// negative case - without token at all
	r, err = http.NewRequest("GET", "test", nil)
	if err != nil {
		t.Error(err)
	}

	res, err = fromHeader(r)
	if err == nil {
		t.Errorf("Result of fromHeader is wrong: expected error, got nil")
	}
	if len(res) != 0 {
		t.Errorf("Expected empty token, got %s", res)
	}

	// negative case - wrong token
	r, err = http.NewRequest("GET", "test", nil)
	if err != nil {
		t.Error(err)
	}
	r.Header.Set("Authorization", "Bearer 123"+token)

	res, err = fromHeader(r)
	if err != nil {
		t.Error(err)
	}

	if res == token {
		t.Error("Token must be wrong, got right")
	}
}

func TestFromBody(t *testing.T) {
	token := "mF_9.B5f-4.1JqM"

	// positive case
	c := new(router.Control)
	data := url.Values{}
	data.Set("access_token", token)
	r, err := http.NewRequest("POST", "test", bytes.NewBufferString(data.Encode()))
	if err != nil {
		t.Error(err)
	}
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Request = r

	res, err := fromBody(r)
	if err != nil {
		t.Error(err)
	}
	if res != token {
		t.Errorf("Token %s is wrong, expected: %s", res, token)
	}

	// negative test, wrong content type
	c = new(router.Control)
	data = url.Values{}
	data.Set("access_token", token)
	r, err = http.NewRequest("POST", "test", bytes.NewBufferString(data.Encode()))
	if err != nil {
		t.Error(err)
	}
	r.Header.Set("Content-Type", "application/json")
	c.Request = r

	res, err = fromBody(r)
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if len(res) != 0 {
		t.Errorf("Expected empty token, got %s", res)
	}

	// negative test, no token at all
	c = new(router.Control)
	data = url.Values{}
	r, err = http.NewRequest("POST", "test", bytes.NewBufferString(data.Encode()))
	if err != nil {
		t.Error(err)
	}
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Request = r

	res, err = fromBody(r)
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if len(res) != 0 {
		t.Errorf("Expected empty token, got %s", res)
	}
}

func TestFromQueryString(t *testing.T) {
	token := "mF_9.B5f-4.1JqM"

	// positive case
	c := new(router.Control)
	r, err := http.NewRequest("GET", "test", nil)
	if err != nil {
		t.Error(err)
	}
	c.Request = r

	param := router.Param{
		Key:   "access_token",
		Value: token,
	}
	c.Set(param)

	res, err := fromQueryString(c)
	if err != nil {
		t.Error(err)
	}
	if res != token {
		t.Errorf("Token %s is wrong, expected: %s", res, token)
	}

	// negative test, no token at all
	c = new(router.Control)
	r, err = http.NewRequest("GET", "test", nil)
	if err != nil {
		t.Error(err)
	}
	c.Request = r

	res, err = fromQueryString(c)
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if len(res) != 0 {
		t.Errorf("Expected empty token, got %s", res)
	}
}

func TestGetBearerTokenFromHeader(t *testing.T) {
	token := "mF_9.B5f-4.1JqM"

	// - Test Set #1. From Header. -
	// positive case
	c := new(router.Control)
	r, err := http.NewRequest("GET", "test", nil)
	if err != nil {
		t.Error(err)
	}
	r.Header.Set("Authorization", "Bearer "+token)
	c.Request = r

	res, err := getBearerToken(c)
	if err != nil {
		t.Error(err)
	}
	if res != token {
		t.Errorf("Token %s is wrong, expected: %s", res, token)
	}

	// negative case - without token at all
	c = new(router.Control)
	r, err = http.NewRequest("GET", "test", nil)
	if err != nil {
		t.Error(err)
	}
	c.Request = r

	res, err = getBearerToken(c)
	if err == nil {
		t.Errorf("Result of fromHeader is wrong: expected error, got nil")
	}
	if len(res) != 0 {
		t.Errorf("Expected empty token, got %s", res)
	}

	// negative case - wrong token
	c = new(router.Control)
	r, err = http.NewRequest("GET", "test", nil)
	if err != nil {
		t.Error(err)
	}
	r.Header.Set("Authorization", "Bearer 123"+token)
	c.Request = r

	res, err = getBearerToken(c)
	if err != nil {
		t.Error(err)
	}

	if res == token {
		t.Error("Token must be wrong, got right")
	}
}

func TestGetBearerTokenFromBody(t *testing.T) {
	token := "mF_9.B5f-4.1JqM"

	// - Test Set #2. From Body. -
	// positive case
	c := new(router.Control)
	data := url.Values{}
	data.Set("access_token", token)
	r, err := http.NewRequest("POST", "test", bytes.NewBufferString(data.Encode()))
	if err != nil {
		t.Error(err)
	}
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Request = r

	res, err := getBearerToken(c)
	if err != nil {
		t.Error(err)
	}
	if res != token {
		t.Errorf("Token %s is wrong, expected: %s", res, token)
	}

	// negative test, wrong content type
	c = new(router.Control)
	data = url.Values{}
	data.Set("access_token", token)
	r, err = http.NewRequest("POST", "test", bytes.NewBufferString(data.Encode()))
	if err != nil {
		t.Error(err)
	}
	r.Header.Set("Content-Type", "application/json")
	c.Request = r

	res, err = getBearerToken(c)
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if len(res) != 0 {
		t.Errorf("Expected empty token, got %s", res)
	}

	// negative test, no token at all
	c = new(router.Control)
	data = url.Values{}
	r, err = http.NewRequest("POST", "test", bytes.NewBufferString(data.Encode()))
	if err != nil {
		t.Error(err)
	}
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	c.Request = r

	res, err = getBearerToken(c)
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if len(res) != 0 {
		t.Errorf("Expected empty token, got %s", res)
	}
}

func TestGetBearerTokenFromQueryString(t *testing.T) {
	token := "mF_9.B5f-4.1JqM"

	// - Test Set #3. From QueryString. -
	// positive case
	c := new(router.Control)
	r, err := http.NewRequest("GET", "test", nil)
	if err != nil {
		t.Error(err)
	}
	c.Request = r

	param := router.Param{
		Key:   "access_token",
		Value: token,
	}
	c.Set(param)

	res, err := getBearerToken(c)
	if err != nil {
		t.Error(err)
	}
	if res != token {
		t.Errorf("Token %s is wrong, expected: %s", res, token)
	}

	// negative test, no token at all
	c = new(router.Control)
	r, err = http.NewRequest("GET", "test", nil)
	if err != nil {
		t.Error(err)
	}
	c.Request = r

	res, err = getBearerToken(c)
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if len(res) != 0 {
		t.Errorf("Expected empty token, got %s", res)
	}
}
