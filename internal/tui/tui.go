package tui

import (
	"fmt"

	"endnet-cli/internal/config"
	"endnet-cli/pkg/models"
)

// Runner exposes the ability to drive an interactive terminal interface.
type Runner interface {
	Run(cfg *config.Config, spec models.EndnetSpec, state *models.RemoteState, plan *models.Plan) error
}

// Model encapsulates the data required by the TUI layer.
type Model struct {
	Config *config.Config
	Spec   models.EndnetSpec
	State  *models.RemoteState
	Plan   *models.Plan
}

// Application is a stub implementation of the TUI runner.
type Application struct{}

// NewRunner constructs a TUI runner instance.
func NewRunner() Runner {
	return &Application{}
}

// Run currently prints a short summary until the real TUI is implemented.
func (a *Application) Run(cfg *config.Config, spec models.EndnetSpec, state *models.RemoteState, plan *models.Plan) error {
	model := Model{Config: cfg, Spec: spec, State: state, Plan: plan}
	fmt.Printf("Launching TUI for project %s at %s\n", model.Config.Project, model.Config.Location)
	fmt.Printf("Desired network: %s (%s)\n", model.Spec.Network.Name, model.Spec.Network.CIDR)
	fmt.Printf("Planned operations: %d\n", len(flattenOperations(plan)))
	return nil
}

func flattenOperations(plan *models.Plan) []models.Operation {
	if plan == nil {
		return nil
	}
	var ops []models.Operation
	for _, group := range [][]models.Operation{plan.NetworkOps, plan.ServerOps, plan.FirewallOps, plan.DNSOps} {
		ops = append(ops, group...)
	}
	return ops
}
