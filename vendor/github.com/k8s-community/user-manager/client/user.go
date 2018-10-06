package client

const (
	userUrlStr = "/sync-user"
)

// UserService defines
type UserService struct {
	client *Client
}

// Sync sends request for user activation in k8s system
func (u *UserService) Sync(user *User) (*Token, *Response, error) {
	req, err := u.client.NewRequest(putMethod, userUrlStr, user)
	if err != nil {
		return nil, nil, err
	}

	token := &Token{}
	resp, err := u.client.Do(req, token)
	if err != nil {
		return nil, resp, err
	}

	return token, resp, nil
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
