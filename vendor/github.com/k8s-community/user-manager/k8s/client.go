package k8s

import (
	"fmt"
	"time"

	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
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
		Host:        baseURL,
		BearerToken: bearerToken,
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
