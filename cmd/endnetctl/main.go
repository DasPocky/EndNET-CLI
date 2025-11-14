package main

import (
	"flag"
	"fmt"
	"log"

	"endnet-cli/internal/config"
	"endnet-cli/internal/state"
	"endnet-cli/internal/tasks"
	"endnet-cli/internal/tui"
	"endnet-cli/pkg/models"
	"endnet-cli/pkg/util"
)

func main() {
	var configPath string
	var useTUI bool
	var planOnly bool

	flag.StringVar(&configPath, "config", "config.yaml", "Path to the EndNET configuration file")
	flag.BoolVar(&useTUI, "tui", false, "Launch the interactive terminal UI")
	flag.BoolVar(&planOnly, "plan", false, "Generate an execution plan without applying it")
	flag.Parse()

	loader := config.NewLoader()
	cfg, err := loader.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	spec := cfg.ToSpec()

	retriever := state.NewRetriever()
	currentState, err := retriever.Current(spec)
	if err != nil {
		log.Fatalf("failed to obtain current state: %v", err)
	}

	planner := tasks.NewPlanner()
	plan, err := planner.Plan(spec, currentState)
	if err != nil {
		log.Fatalf("failed to generate plan: %v", err)
	}

	if useTUI {
		runner := tui.NewRunner()
		if err := runner.Run(cfg, spec, currentState, plan); err != nil {
			log.Fatalf("tui exited with error: %v", err)
		}
		return
	}

	if planOnly {
		fmt.Println("Planned actions:")
		for _, op := range flattenPlan(plan) {
			fmt.Printf("- [%s] %s -> %s\n", op.Type, op.Target, op.Details)
		}
		return
	}

	executor := tasks.NewExecutor(util.NewLogger())
	result, err := executor.Execute(plan)
	if err != nil {
		log.Fatalf("plan execution failed: %v", err)
	}

	if result.ChangesApplied {
		fmt.Println("Plan executed successfully.")
	} else {
		fmt.Println("Plan execution completed without changes.")
	}
}

func flattenPlan(plan *models.Plan) []models.Operation {
	if plan == nil {
		return nil
	}
	var ops []models.Operation
	for _, group := range [][]models.Operation{plan.NetworkOps, plan.ServerOps, plan.FirewallOps, plan.DNSOps} {
		ops = append(ops, group...)
	}
	return ops
}
