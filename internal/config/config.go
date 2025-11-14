package config

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"endnet-cli/pkg/models"
)

// Config captures all configuration knobs for EndNET-CLI.
type Config struct {
	Project  string        `yaml:"project"`
	Location string        `yaml:"location"`
	Network  NetworkConfig `yaml:"network"`
	Roles    RolesConfig   `yaml:"roles"`
	DNS      DNSConfig     `yaml:"dns"`
	Hetzner  HetznerConfig `yaml:"hetzner"`
	IPv64    IPv64Config   `yaml:"ipv64"`
	LoadedAt time.Time     `yaml:"-"`
	Source   string        `yaml:"-"`
}

// NetworkConfig contains network defaults.
type NetworkConfig struct {
	Name       string `yaml:"name"`
	CIDR       string `yaml:"cidr"`
	SubnetCIDR string `yaml:"subnetCidr"`
	GatewayIP  string `yaml:"gatewayIp"`
}

// RolesConfig groups all server roles.
type RolesConfig struct {
	Edge  NodeConfig `yaml:"edge"`
	WG    NodeConfig `yaml:"wg"`
	Forge NodeConfig `yaml:"forge"`
}

// NodeConfig defines per-node options.
type NodeConfig struct {
	Name        string `yaml:"name"`
	Type        string `yaml:"type"`
	Image       string `yaml:"image"`
	PrivateIP   string `yaml:"privateIp"`
	HasPublicIP bool   `yaml:"publicIp"`
}

// DNSConfig contains DNS integration options.
type DNSConfig struct {
	RootDomain  string `yaml:"rootDomain"`
	ForgejoHost string `yaml:"forgejoHost"`
}

// HetznerConfig provides credentials and options for Hetzner Cloud.
type HetznerConfig struct {
	APIToken   string `yaml:"apiToken"`
	SSHKeyName string `yaml:"sshKeyName"`
}

// IPv64Config describes the IPv64 API credentials.
type IPv64Config struct {
	APIKey      string `yaml:"apiKey"`
	DynDNSToken string `yaml:"dynDnsToken"`
}

// Loader defines the interface for materialising configuration data.
type Loader interface {
	Load(path string) (*Config, error)
}

// FileLoader implements Loader via YAML files and environment overrides.
type FileLoader struct{}

type yamlEntry struct {
	level int
	key   string
}

// NewLoader instantiates a FileLoader.
func NewLoader() Loader {
	return &FileLoader{}
}

// DefaultConfig returns a fully-populated configuration with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		Project:  "endnet",
		Location: "nbg1",
		Network: NetworkConfig{
			Name:       "endnet-internal",
			CIDR:       "10.10.0.0/16",
			SubnetCIDR: "10.10.0.0/24",
			GatewayIP:  "10.10.0.2",
		},
		Roles: RolesConfig{
			Edge: NodeConfig{
				Name:        "endnet-edge-1",
				Type:        "cx23",
				Image:       "debian-12",
				PrivateIP:   "10.10.0.2",
				HasPublicIP: true,
			},
			WG: NodeConfig{
				Name:        "endnet-wg-1",
				Type:        "cx23",
				Image:       "debian-12",
				PrivateIP:   "10.10.0.10",
				HasPublicIP: false,
			},
			Forge: NodeConfig{
				Name:        "endnet-git-1",
				Type:        "cx23",
				Image:       "debian-12",
				PrivateIP:   "10.10.0.20",
				HasPublicIP: false,
			},
		},
		DNS: DNSConfig{
			RootDomain:  "endnet.ipv64.net",
			ForgejoHost: "git.endnet.ipv64.net",
		},
		Hetzner: HetznerConfig{
			SSHKeyName: "endnet",
		},
		IPv64: IPv64Config{},
	}
}

// Load reads configuration from disk, merging it with defaults and environment overrides.
func (l *FileLoader) Load(path string) (*Config, error) {
	if path == "" {
		return nil, errors.New("configuration path must not be empty")
	}

	cfg := DefaultConfig()
	cfg.Source = path

	if err := l.applyFile(path, cfg); err != nil {
		return nil, err
	}

	l.applyEnv(cfg)
	cfg.LoadedAt = time.Now()

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (l *FileLoader) applyFile(path string, cfg *Config) error {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("read config: %w", err)
	}

	if err := parseMinimalYAML(data, cfg); err != nil {
		return fmt.Errorf("parse config: %w", err)
	}

	return nil
}

func (l *FileLoader) applyEnv(cfg *Config) {
	if v := os.Getenv("ENDNET_PROJECT"); v != "" {
		cfg.Project = v
	}
	if v := os.Getenv("ENDNET_LOCATION"); v != "" {
		cfg.Location = v
	}
	if v := os.Getenv("ENDNET_HCLOUD_TOKEN"); v != "" {
		cfg.Hetzner.APIToken = v
	}
	if v := os.Getenv("ENDNET_IPV64_API_KEY"); v != "" {
		cfg.IPv64.APIKey = v
	}
	if v := os.Getenv("ENDNET_IPV64_DYNDNS_TOKEN"); v != "" {
		cfg.IPv64.DynDNSToken = v
	}
}

// ToSpec converts the configuration into the desired-state specification.
func (c *Config) ToSpec() models.EndnetSpec {
	extras := make(map[string]models.NodeSpec)

	return models.EndnetSpec{
		Project:  c.Project,
		Location: c.Location,
		Network: models.NetworkSpec{
			Name:       c.Network.Name,
			CIDR:       c.Network.CIDR,
			SubnetCIDR: c.Network.SubnetCIDR,
			GatewayIP:  c.Network.GatewayIP,
		},
		Roles: models.RolesSpec{
			Edge:   toNodeSpec(c.Roles.Edge),
			WG:     toNodeSpec(c.Roles.WG),
			Forge:  toNodeSpec(c.Roles.Forge),
			Extras: extras,
		},
		DNS: models.DNSSpec{
			RootDomain:  c.DNS.RootDomain,
			ForgejoHost: c.DNS.ForgejoHost,
		},
	}
}

func toNodeSpec(cfg NodeConfig) models.NodeSpec {
	return models.NodeSpec{
		Name:        cfg.Name,
		Type:        cfg.Type,
		Image:       cfg.Image,
		PrivateIP:   cfg.PrivateIP,
		HasPublicIP: cfg.HasPublicIP,
	}
}

// Validate performs basic sanity checks on the configuration.
func (c *Config) Validate() error {
	if c.Project == "" {
		return errors.New("project must not be empty")
	}
	if c.Network.Name == "" {
		return errors.New("network.name must not be empty")
	}
	if c.Roles.Edge.Name == "" {
		return errors.New("roles.edge.name must not be empty")
	}
	if c.DNS.RootDomain == "" {
		return errors.New("dns.rootDomain must not be empty")
	}
	return nil
}

func parseMinimalYAML(data []byte, cfg *Config) error {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	var stack []yamlEntry

	setValue := func(path []string, value string) error {
		key := strings.Join(path, ".")
		switch key {
		case "project":
			cfg.Project = value
		case "location":
			cfg.Location = value
		case "network.name":
			cfg.Network.Name = value
		case "network.cidr":
			cfg.Network.CIDR = value
		case "network.subnetCidr":
			cfg.Network.SubnetCIDR = value
		case "network.gatewayIp":
			cfg.Network.GatewayIP = value
		case "roles.edge.name":
			cfg.Roles.Edge.Name = value
		case "roles.edge.type":
			cfg.Roles.Edge.Type = value
		case "roles.edge.image":
			cfg.Roles.Edge.Image = value
		case "roles.edge.privateIp":
			cfg.Roles.Edge.PrivateIP = value
		case "roles.edge.publicIp":
			cfg.Roles.Edge.HasPublicIP = parseBool(value)
		case "roles.wg.name":
			cfg.Roles.WG.Name = value
		case "roles.wg.type":
			cfg.Roles.WG.Type = value
		case "roles.wg.image":
			cfg.Roles.WG.Image = value
		case "roles.wg.privateIp":
			cfg.Roles.WG.PrivateIP = value
		case "roles.wg.publicIp":
			cfg.Roles.WG.HasPublicIP = parseBool(value)
		case "roles.forge.name":
			cfg.Roles.Forge.Name = value
		case "roles.forge.type":
			cfg.Roles.Forge.Type = value
		case "roles.forge.image":
			cfg.Roles.Forge.Image = value
		case "roles.forge.privateIp":
			cfg.Roles.Forge.PrivateIP = value
		case "roles.forge.publicIp":
			cfg.Roles.Forge.HasPublicIP = parseBool(value)
		case "dns.rootDomain":
			cfg.DNS.RootDomain = value
		case "dns.forgejoHost":
			cfg.DNS.ForgejoHost = value
		case "hetzner.apiToken":
			cfg.Hetzner.APIToken = value
		case "hetzner.sshKeyName":
			cfg.Hetzner.SSHKeyName = value
		case "ipv64.apiKey":
			cfg.IPv64.APIKey = value
		case "ipv64.dynDnsToken":
			cfg.IPv64.DynDNSToken = value
		default:
			// ignore unknown keys for now
		}
		return nil
	}

	for scanner.Scan() {
		raw := scanner.Text()
		trimmed := strings.TrimSpace(raw)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		indent := len(raw) - len(strings.TrimLeft(raw, " "))
		level := indent / 2

		for len(stack) > 0 && stack[len(stack)-1].level >= level {
			stack = stack[:len(stack)-1]
		}

		parts := strings.SplitN(trimmed, ":", 2)
		key := strings.TrimSpace(parts[0])
		value := ""
		if len(parts) == 2 {
			value = strings.TrimSpace(parts[1])
		}

		path := append([]string{}, keysFromStack(stack)...)
		path = append(path, key)

		if value == "" {
			stack = append(stack, yamlEntry{level: level, key: key})
			continue
		}

		value = strings.Trim(value, "\"'")
		if err := setValue(path, value); err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func keysFromStack(stack []yamlEntry) []string {
	keys := make([]string, len(stack))
	for i, e := range stack {
		keys[i] = e.key
	}
	return keys
}

func parseBool(value string) bool {
	b, err := strconv.ParseBool(value)
	if err != nil {
		return false
	}
	return b
}
