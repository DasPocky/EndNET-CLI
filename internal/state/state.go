package state

import (
	"fmt"
	"time"

	"endnet-cli/pkg/models"
)

// Retriever gathers state from infrastructure providers.
type Retriever interface {
	Current(cfg *models.Config) (*models.State, error)
}

// Snapshotter is a trivial implementation of Retriever used for bootstrapping.
type Snapshotter struct{}

// NewRetriever creates a Retriever implementation suitable for development
// environments until real integrations are available.
func NewRetriever() Retriever {
	return &Snapshotter{}
}

// Current provides a static snapshot of the world for now.
func (r *Snapshotter) Current(cfg *models.Config) (*models.State, error) {
	if cfg == nil {
		return nil, fmt.Errorf("configuration must not be nil")
	}

	return &models.State{
		RetrievedAt: time.Now(),
		Summary:     "static snapshot",
	}, nil
}
