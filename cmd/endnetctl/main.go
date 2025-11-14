package main

import (
	"flag"
	"fmt"
	"log"

	"endnet-cli/internal/config"
	"endnet-cli/internal/state"
	"endnet-cli/internal/tasks"
	"endnet-cli/internal/tui"
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

	retriever := state.NewRetriever()
	currentState, err := retriever.Current(cfg)
	if err != nil {
		log.Fatalf("failed to obtain current state: %v", err)
	}

	if useTUI {
		runner := tui.NewRunner()
		if err := runner.Run(cfg, currentState); err != nil {
			log.Fatalf("tui exited with error: %v", err)
		}
		return
	}

	planner := tasks.NewPlanner()
	plan, err := planner.Plan(cfg, currentState)
	if err != nil {
		log.Fatalf("failed to generate plan: %v", err)
	}

	if planOnly {
		fmt.Println("Planned actions:")
		for _, step := range plan.Steps {
			fmt.Printf("- %s\n", step)
		}
		return
	}

	executor := tasks.NewExecutor()
	result, err := executor.Execute(plan)
	if err != nil {
		log.Fatalf("plan execution failed: %v", err)
	}

	if result.Applied {
		fmt.Println("Plan executed successfully.")
		return
	}

	fmt.Println("Plan execution completed without changes.")
}
