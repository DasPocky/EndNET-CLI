package state

import (
	"errors"
	"time"

	"endnet-cli/pkg/models"
)

// Retriever gathers state from infrastructure providers.
type Retriever interface {
	Current(spec models.EndnetSpec) (*models.RemoteState, error)
}

// Snapshotter is a placeholder Retriever that returns empty provider state.
type Snapshotter struct {
	Now func() time.Time
}

// NewRetriever constructs a development-friendly Retriever implementation.
func NewRetriever() Retriever {
	return &Snapshotter{Now: time.Now}
}

// Current produces a deterministic RemoteState snapshot for bootstrapping.
func (r *Snapshotter) Current(spec models.EndnetSpec) (*models.RemoteState, error) {
	if spec.Project == "" {
		return nil, errors.New("spec project must not be empty")
	}

	now := r.Now()

	return &models.RemoteState{
		Hetzner: models.HetznerState{
			Networks:  []models.Network{},
			Routes:    []models.Route{},
			Servers:   []models.Server{},
			Firewalls: []models.Firewall{},
		},
		IPv64: models.IPv64State{
			Domains: map[string]models.Domain{},
		},
		RetrievedAt: now,
	}, nil
}
