package hetzner

import (
	"fmt"

	"endnet-cli/pkg/models"
)

// Client describes the operations required to interact with the Hetzner Cloud.
type Client interface {
	Authenticate(token string) error
	ListServers() ([]models.Server, error)
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
	if token == "" {
		return fmt.Errorf("token must not be empty")
	}
	c.token = token
	return nil
}

// ListServers returns a deterministic list so higher layers can be developed.
func (c *APIClient) ListServers() ([]models.Server, error) {
	if c.token == "" {
		return nil, models.ErrUnauthenticated
	}

	return []models.Server{{
		ID:        1,
		Name:      "placeholder-server",
		Type:      "cx23",
		Image:     "debian-12",
		PrivateIP: "10.10.0.2",
	}}, nil
}
