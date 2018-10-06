package k8s

import (
	"fmt"
	"strings"
	"time"

	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/apis/rbac/v1beta1"
	"k8s.io/client-go/rest"
)

// Client defines
type Client struct {
	// HTTP client used to communicate with the API.
	client *kubernetes.Clientset

	// Base URL for API requests.
	BaseURL string

	// Services used for talking to different parts of the API.
	BearerToken string
}

// NewClient initializes client for k8s API
func NewClient(baseURL string, bearerToken string) (*Client, error) {
	config := &rest.Config{
		Host:            baseURL,
		BearerToken:     bearerToken,
		TLSClientConfig: rest.TLSClientConfig{Insecure: true}, // todo: use cacert instead
	}

	c, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("cannot connect to k8s server: %s", err)
	}

	return &Client{
		client:      c,
		BaseURL:     baseURL,
		BearerToken: bearerToken,
	}, nil
}

// GetNamespace creates namespace with defined name in k8s system
func (c *Client) GetNamespace(namespaceName string) (*v1.Namespace, error) {
	options := meta_v1.GetOptions{}
	namespace, err := c.client.Core().Namespaces().Get(namespaceName, options)
	if err != nil {
		return nil, fmt.Errorf("cannot get namespace: %s", err)
	}

	return namespace, nil
}

// CreateNamespace creates namespace with defined name in k8s system
func (c *Client) CreateNamespace(namespaceName string) error {
	namespace := &v1.Namespace{}
	namespace.APIVersion = "v1"
	namespace.Kind = "Namespace"
	namespace.ObjectMeta.Name = namespaceName
	namespace.ObjectMeta.Labels = make(map[string]string)
	namespace.ObjectMeta.Labels["workshop/date"] = time.Now().Format("2006-01-02")

	_, err := c.client.Core().Namespaces().Create(namespace)
	if err != nil {
		return fmt.Errorf("cannot create namespace: %s", err)
	}

	return nil
}

// CreateNamespaceAdmin creates admin service account.
func (c *Client) CreateNamespaceAdmin(namespace string) error {
	sa := &v1.ServiceAccount{}
	sa.APIVersion = "v1"
	sa.Kind = "ServiceAccount"
	sa.ObjectMeta.Name = namespace
	sa.ObjectMeta.Namespace = namespace

	_, err := c.client.ServiceAccounts(namespace).Create(sa)
	if err != nil {
		return fmt.Errorf("cannot create service account: %s", err)
	}

	rb := &v1beta1.RoleBinding{}
	subj := v1beta1.Subject{
		Kind:      "ServiceAccount",
		Name:      namespace,
		Namespace: namespace,
	}
	rb.Kind = "RoleBinding"
	rb.ObjectMeta.Name = namespace + "-admin"
	rb.Subjects = append(rb.Subjects, subj)
	rb.RoleRef = v1beta1.RoleRef{
		APIGroup: "rbac.authorization.k8s.io",
		Name:     "admin",
		Kind:     "ClusterRole",
	}

	_, err = c.client.RbacV1beta1Client.RoleBindings(namespace).Create(rb)
	if err != nil {
		return fmt.Errorf("cannot create role binding: %s", err)
	}

	return nil
}

func (c *Client) GetNamespaceToken(namespace string) (*Token, error) {
	secrets, err := c.client.Secrets(namespace).List(meta_v1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("cannot get lis of secrets: %s", err)
	}

	var tokenSecr v1.Secret
	found := false
	for _, secr := range secrets.Items {
		if strings.HasPrefix(secr.ObjectMeta.Name, namespace+"-token") {
			tokenSecr = secr
			found = true
			break
		}
	}

	if !found {
		return nil, fmt.Errorf("cannot find token secret")
	}

	tok := Token{
		Cert:  string(tokenSecr.Data["ca.crt"]),
		Token: string(tokenSecr.Data["token"]),
	}
	return &tok, nil
}

// CopySecret creates copy of secret with name secretName from defined namespace
// to other namespace in k8s system
func (c *Client) CopySecret(secretName string, fromNamespace string, toNamespace string) error {
	options := meta_v1.GetOptions{}
	secret, err := c.client.Core().Secrets(fromNamespace).Get(secretName, options)
	if err != nil {
		return fmt.Errorf("cannot get secret %s from namespace %s: %s", secretName, fromNamespace, err)
	}

	newSecret := &v1.Secret{}
	newSecret.APIVersion = secret.APIVersion
	newSecret.Name = secret.Name
	newSecret.Kind = secret.Kind
	newSecret.Labels = secret.Labels
	newSecret.Annotations = secret.Annotations
	newSecret.Data = secret.Data
	newSecret.Type = secret.Type

	newSecret.Namespace = toNamespace

	newSecret, err = c.client.Core().Secrets(toNamespace).Create(newSecret)
	if err != nil {
		return fmt.Errorf("cannot create secret %s in namespace %s: %s", secretName, toNamespace, err)
	}

	return nil
}
