package tasks

import (
	"fmt"

	"endnet-cli/internal/cloudinit"
	"endnet-cli/pkg/models"
)

// Planner builds plans for reconciling desired state with the current world.
type Planner interface {
	Plan(cfg *models.Config, state *models.State) (*models.Plan, error)
}

// Executor applies plans and reports their results.
type Executor interface {
	Execute(plan *models.Plan) (*models.Result, error)
}

// Pipeline wires together the planner, executor, and any helper services.
type Pipeline struct {
	generator cloudinit.Generator
}

// NewPlanner returns a Planner backed by the Pipeline placeholder implementation.
func NewPlanner() Planner {
	return &Pipeline{generator: cloudinit.NewGenerator()}
}

// NewExecutor returns an Executor backed by the Pipeline placeholder implementation.
func NewExecutor() Executor {
	return &Pipeline{generator: cloudinit.NewGenerator()}
}

// Plan produces a deterministic sequence of steps while more advanced logic is in development.
func (p *Pipeline) Plan(cfg *models.Config, state *models.State) (*models.Plan, error) {
	if cfg == nil || state == nil {
		return nil, fmt.Errorf("both configuration and state are required")
	}

	userData, err := p.generator.Generate(cfg, &models.Plan{})
	if err != nil {
		return nil, err
	}

	return &models.Plan{
		Steps: []string{
			"Validate configuration",
			"Synchronize infrastructure state",
			fmt.Sprintf("Render cloud-init user-data: %s", userData),
		},
	}, nil
}

// Execute returns a successful result without applying any real changes yet.
func (p *Pipeline) Execute(plan *models.Plan) (*models.Result, error) {
	if plan == nil {
		return nil, fmt.Errorf("plan must not be nil")
	}

	return &models.Result{Applied: len(plan.Steps) > 0}, nil
}
