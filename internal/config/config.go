package config

import (
	"fmt"

	"endnet-cli/pkg/models"
)

// Loader knows how to materialize configuration data for EndNET.
type Loader interface {
	Load(path string) (*models.Config, error)
}

// FileLoader implements Loader by reading configuration files from disk.
type FileLoader struct{}

// NewLoader returns a Loader that can read EndNET configuration files.
func NewLoader() Loader {
	return &FileLoader{}
}

// Load returns a placeholder configuration structure so that other
// components can be wired together while real loading logic is developed.
func (l *FileLoader) Load(path string) (*models.Config, error) {
	if path == "" {
		return nil, fmt.Errorf("configuration path must not be empty")
	}

	return &models.Config{
		Source: path,
		Metadata: map[string]string{
			"environment": "development",
		},
	}, nil
}
