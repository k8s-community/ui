package bearer

import (
	"errors"
	"net/http"
	"strings"

	"github.com/takama/router"
)

// WithToken checks if token from request suits the given token.
func WithToken(expectedToken string, handler router.Handle) router.Handle {
	return func(ctx *router.Control) {
		token, err := getBearerToken(ctx)
		if err != nil {
			http.Error(ctx.Writer, err.Error(), http.StatusBadRequest)
			return
		}

		if token != expectedToken {
			http.Error(ctx.Writer, "Bearer token is incorrect", http.StatusBadRequest)
			return
		}

		handler(ctx)
	}
}

// getBearerToken try to find access token in header, body and query string.
func getBearerToken(ctx *router.Control) (string, error) {
	token, err := fromHeader(ctx.Request)
	if err == nil {
		return token, nil
	}

	token, err = fromBody(ctx.Request)
	if err == nil {
		return token, nil
	}

	token, err = fromQueryString(ctx)
	if err == nil {
		return token, nil
	}

	return "", errors.New("Could not get an access token from the request")
}

// fromHeader parse auth token from request header: https://tools.ietf.org/html/rfc6750#section-2.1
func fromHeader(request *http.Request) (string, error) {
	authHeaderKey := "Authorization"
	authHeaderValue := request.Header.Get(authHeaderKey)
	authScheme := "Bearer"

	l := len(authScheme)
	if len(authHeaderValue) > l+1 && authHeaderValue[:l] == authScheme {
		return authHeaderValue[l+1:], nil
	}

	return "", errors.New("The request header does not contain an access token")
}

// fromBody parse auth token from form: https://tools.ietf.org/html/rfc6750#section-2.2
func fromBody(request *http.Request) (string, error) {
	contentType := request.Header.Get("Content-Type")

	if strings.ToLower(contentType) != "application/x-www-form-urlencoded" {
		return "", errors.New("Body access token is not supported by request content type")
	}

	authParamName := "access_token"
	authParamValue := request.FormValue(authParamName)

	if len(authParamValue) > 0 {
		return authParamValue, nil
	}

	return "", errors.New("The request body does not contain an access token")
}

// fromQueryString parse auth token from query string: https://tools.ietf.org/html/rfc6750#section-2.3
func fromQueryString(ctx *router.Control) (string, error) {
	authParamName := "access_token"
	authParamValue := ctx.Get(authParamName)

	if len(authParamValue) > 0 {
		return authParamValue, nil
	}

	return "", errors.New("The query string does not contain an access token")
}
