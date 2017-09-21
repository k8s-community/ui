package client

import "net/http"

const (
	userUrlStr = "/sync-user"
)

// UserService defines
type UserService struct {
	client *Client
}

// Sync sends request for user activation in k8s system
func (u *UserService) Sync(user *User) (int, error) {
	req, err := u.client.NewRequest(putMethod, userUrlStr, user)
	if err != nil {
		return http.StatusBadRequest, err
	}

	resp, err := u.client.Do(req, nil)
	if err != nil {
		return http.StatusBadRequest, err
	}

	return resp.StatusCode, nil
}

// User defines
type User struct {
	Name string `json:"name"`
}

// NewUser create a new User instance
func NewUser(username string) *User {
	return &User{
		Name: username,
	}
}
