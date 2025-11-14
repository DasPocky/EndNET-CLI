package tui

import (
	"fmt"

	"endnet-cli/pkg/models"
)

// Runner exposes the ability to drive an interactive terminal interface.
type Runner interface {
	Run(cfg *models.Config, state *models.State) error
}

// Model encapsulates the data required by the TUI layer.
type Model struct {
	Config *models.Config
	State  *models.State
}

// Application is a stub implementation of the TUI runner.
type Application struct{}

// NewRunner constructs a TUI runner instance.
func NewRunner() Runner {
	return &Application{}
}

// Run currently prints a short summary until the real TUI is implemented.
func (a *Application) Run(cfg *models.Config, state *models.State) error {
	model := Model{Config: cfg, State: state}
	fmt.Printf("Launching TUI with config %s and state summary %s\n", model.Config.Source, model.State.Summary)
	return nil
}
