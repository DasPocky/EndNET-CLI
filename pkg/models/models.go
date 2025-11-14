package models

import "errors"

// Config represents the desired state configuration.
type Config struct {
	Source   string
	Metadata map[string]string
}

// State captures the observed state from providers.
type State struct {
	Summary     string
	RetrievedAt interface{}
}

// Plan enumerates the actions required to converge state.
type Plan struct {
	Steps []string
}

// Result indicates whether the reconciliation run applied any changes.
type Result struct {
	Applied bool
}

// HetznerServer is a placeholder model describing a provisioned server.
type HetznerServer struct {
	Name string
}

// DNSRecord is a placeholder DNS entry representation.
type DNSRecord struct {
	Name    string
	Address string
}

// ErrUnauthenticated is returned when an API client is used before authentication.
var ErrUnauthenticated = errors.New("client is not authenticated")
