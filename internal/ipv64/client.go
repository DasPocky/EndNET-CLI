package ipv64

import "endnet-cli/pkg/models"

// Client exposes the subset of the IPv64 API required by the controller.
type Client interface {
	Authenticate(token string) error
	ListRecords() ([]models.DNSRecord, error)
}

// APIClient is a stubbed IPv64 implementation.
type APIClient struct {
	token string
}

// NewClient returns a placeholder IPv64 client implementation.
func NewClient() Client {
	return &APIClient{}
}

// Authenticate stores the API token for later use.
func (c *APIClient) Authenticate(token string) error {
	c.token = token
	return nil
}

// ListRecords returns a static set of DNS records for now.
func (c *APIClient) ListRecords() ([]models.DNSRecord, error) {
	if c.token == "" {
		return nil, models.ErrUnauthenticated
	}

	return []models.DNSRecord{{Name: "example.com", Address: "192.0.2.1"}}, nil
}
