package k8s

// Token represents token type
type Token struct {
	Cert  string `json:"cert"`
	Token string `json:"token"`
}
