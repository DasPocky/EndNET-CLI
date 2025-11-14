package hetzner

import "endnet-cli/pkg/models"

// Client describes the operations required to interact with the Hetzner Cloud.
type Client interface {
	Authenticate(token string) error
	ListServers() ([]models.HetznerServer, error)
}

// APIClient is a stub implementation until the real SDK integration is written.
type APIClient struct {
	token string
}

// NewClient returns a placeholder Hetzner client.
func NewClient() Client {
	return &APIClient{}
}

// Authenticate stores the provided token for later use.
func (c *APIClient) Authenticate(token string) error {
	c.token = token
	return nil
}

// ListServers returns a deterministic list so higher layers can be developed.
func (c *APIClient) ListServers() ([]models.HetznerServer, error) {
	if c.token == "" {
		return nil, ErrUnauthenticated
	}

	return []models.HetznerServer{{Name: "placeholder-server"}}, nil
}

// ErrUnauthenticated is returned when the client has not been initialized with a token.
var ErrUnauthenticated = models.ErrUnauthenticated
