package tasks

import (
	"fmt"

	"endnet-cli/pkg/models"
)

// Planner builds plans for reconciling desired and current state.
type Planner interface {
	Plan(spec models.EndnetSpec, state *models.RemoteState) (*models.Plan, error)
}

// DefaultPlanner is a basic planner implementation that focuses on
// producing human-readable operations for the TUI and CLI output.
type DefaultPlanner struct{}

// NewPlanner returns a planner ready for use.
func NewPlanner() Planner {
	return &DefaultPlanner{}
}

// Plan compares the desired specification with the observed state and
// generates a list of operations required to reconcile them.
func (p *DefaultPlanner) Plan(spec models.EndnetSpec, state *models.RemoteState) (*models.Plan, error) {
	if state == nil {
		return nil, fmt.Errorf("state must not be nil")
	}

	plan := &models.Plan{}

	if !hasNetwork(state.Hetzner.Networks, spec.Network.Name) {
		plan.NetworkOps = append(plan.NetworkOps, models.Operation{
			Type:    "create",
			Target:  fmt.Sprintf("network:%s", spec.Network.Name),
			Details: fmt.Sprintf("create network %s with %s", spec.Network.Name, spec.Network.CIDR),
		})
	} else {
		plan.NetworkOps = append(plan.NetworkOps, models.Operation{
			Type:    "noop",
			Target:  fmt.Sprintf("network:%s", spec.Network.Name),
			Details: "network already present",
		})
	}

	ensureServer(plan, state, spec.Roles.Edge)
	ensureServer(plan, state, spec.Roles.WG)
	ensureServer(plan, state, spec.Roles.Forge)

	if _, ok := state.IPv64.Domains[spec.DNS.RootDomain]; !ok {
		plan.DNSOps = append(plan.DNSOps, models.Operation{
			Type:    "verify",
			Target:  fmt.Sprintf("dns:%s", spec.DNS.RootDomain),
			Details: "ensure root domain exists in IPv64 account",
		})
	} else {
		plan.DNSOps = append(plan.DNSOps, models.Operation{
			Type:    "noop",
			Target:  fmt.Sprintf("dns:%s", spec.DNS.RootDomain),
			Details: "root domain present",
		})
	}

	plan.DNSOps = append(plan.DNSOps, models.Operation{
		Type:    "update",
		Target:  fmt.Sprintf("dns:A %s", spec.DNS.RootDomain),
		Details: "synchronize A record using DynDNS",
	})

	plan.DNSOps = append(plan.DNSOps, models.Operation{
		Type:    "ensure",
		Target:  fmt.Sprintf("dns:CNAME %s", spec.DNS.ForgejoHost),
		Details: fmt.Sprintf("ensure CNAME points to %s", spec.DNS.RootDomain),
	})

	plan.FirewallOps = append(plan.FirewallOps, models.Operation{
		Type:    "reconcile",
		Target:  "firewall:endnet-edge",
		Details: "ensure firewall rules for edge host",
	})

	return plan, nil
}

func ensureServer(plan *models.Plan, state *models.RemoteState, node models.NodeSpec) {
	if node.Name == "" {
		return
	}

	if !hasServer(state.Hetzner.Servers, node.Name) {
		plan.ServerOps = append(plan.ServerOps, models.Operation{
			Type:    "create",
			Target:  fmt.Sprintf("server:%s", node.Name),
			Details: fmt.Sprintf("provision %s (%s) with IP %s", node.Type, node.Image, node.PrivateIP),
		})
	} else {
		plan.ServerOps = append(plan.ServerOps, models.Operation{
			Type:    "noop",
			Target:  fmt.Sprintf("server:%s", node.Name),
			Details: "server already exists",
		})
	}
}

func hasNetwork(networks []models.Network, name string) bool {
	for _, n := range networks {
		if n.Name == name {
			return true
		}
	}
	return false
}

func hasServer(servers []models.Server, name string) bool {
	for _, s := range servers {
		if s.Name == name {
			return true
		}
	}
	return false
}
