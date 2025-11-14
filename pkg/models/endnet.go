package models

import (
	"errors"
	"time"
)

// EndnetSpec represents the desired state of the EndNET infrastructure.
type EndnetSpec struct {
	Project  string
	Location string
	Network  NetworkSpec
	Roles    RolesSpec
	DNS      DNSSpec
}

// NetworkSpec contains the required network configuration.
type NetworkSpec struct {
	Name       string
	CIDR       string
	SubnetCIDR string
	GatewayIP  string
}

// RolesSpec enumerates the infrastructure roles that should exist.
type RolesSpec struct {
	Edge   NodeSpec
	WG     NodeSpec
	Forge  NodeSpec
	Extras map[string]NodeSpec
}

// NodeSpec describes a single server instance.
type NodeSpec struct {
	Name        string
	Type        string
	Image       string
	PrivateIP   string
	HasPublicIP bool
}

// DNSSpec details the DNS records required for the infrastructure.
type DNSSpec struct {
	RootDomain  string
	ForgejoHost string
}

// RemoteState captures the current view of the providers.
type RemoteState struct {
	Hetzner     HetznerState
	IPv64       IPv64State
	RetrievedAt time.Time
}

// HetznerState mirrors the Hetzner resources relevant to EndNET.
type HetznerState struct {
	Networks  []Network
	Routes    []Route
	Servers   []Server
	Firewalls []Firewall
}

// Network represents a Hetzner network.
type Network struct {
	ID   int
	Name string
	CIDR string
}

// Route represents a network route.
type Route struct {
	NetworkID       int
	DestinationCIDR string
	GatewayIP       string
}

// Server describes a provisioned Hetzner server.
type Server struct {
	ID        int
	Name      string
	Type      string
	Image     string
	PrivateIP string
	PublicIP  string
}

// Firewall captures firewall configuration details.
type Firewall struct {
	ID    int
	Name  string
	Rules []FirewallRule
}

// FirewallRule describes a single firewall rule.
type FirewallRule struct {
	Direction string
	Protocol  string
	Port      string
	Source    string
	Target    string
}

// IPv64State holds DNS domain information.
type IPv64State struct {
	Domains map[string]Domain
}

// Domain describes DNS records for a specific domain.
type Domain struct {
	Name    string
	Records []DNSRecord
}

// DNSRecord represents a DNS record entry.
type DNSRecord struct {
	Type  string
	Name  string
	Value string
	TTL   int
}

// Plan summarizes the operations necessary to reach the desired state.
type Plan struct {
	NetworkOps  []Operation
	ServerOps   []Operation
	FirewallOps []Operation
	DNSOps      []Operation
}

// Operation is a single action in a plan.
type Operation struct {
	Type    string
	Target  string
	Details string
}

// ExecutionResult captures the outcome of applying a plan.
type ExecutionResult struct {
	ChangesApplied    bool
	AppliedOperations []Operation
	StartedAt         time.Time
	CompletedAt       time.Time
	Notes             []string
}

// ErrUnauthenticated indicates API usage prior to authentication.
var ErrUnauthenticated = errors.New("client is not authenticated")
