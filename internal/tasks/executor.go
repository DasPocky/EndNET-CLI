package tasks

import (
	"errors"
	"time"

	"endnet-cli/pkg/models"
	"endnet-cli/pkg/util"
)

// Executor applies plans and reports their results.
type Executor interface {
	Execute(plan *models.Plan) (*models.ExecutionResult, error)
}

// DefaultExecutor logs the operations that would be executed.
type DefaultExecutor struct {
	logger util.Logger
}

// NewExecutor constructs an Executor with the provided logger.
func NewExecutor(logger util.Logger) Executor {
	if logger == nil {
		logger = util.NewLogger()
	}
	return &DefaultExecutor{logger: logger}
}

// Execute iterates through plan operations and logs them.
func (e *DefaultExecutor) Execute(plan *models.Plan) (*models.ExecutionResult, error) {
	if plan == nil {
		return nil, errors.New("plan must not be nil")
	}

	started := time.Now()
	var applied []models.Operation

	for _, ops := range [][]models.Operation{plan.NetworkOps, plan.ServerOps, plan.FirewallOps, plan.DNSOps} {
		for _, op := range ops {
			e.logger.Infof("%s %s (%s)", op.Type, op.Target, op.Details)
			applied = append(applied, op)
		}
	}

	return &models.ExecutionResult{
		ChangesApplied:    len(applied) > 0,
		AppliedOperations: applied,
		StartedAt:         started,
		CompletedAt:       time.Now(),
		Notes:             []string{"dry-run execution"},
	}, nil
}
