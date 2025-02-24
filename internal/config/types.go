package config

// Config represents the complete iptables configuration
type Config struct {
	Tables map[string]*TableConfig `yaml:"tables"`
}

// TableConfig represents a single iptables table configuration
type TableConfig struct {
	Chains map[string]*ChainConfig `yaml:"chains"`
}

// ChainConfig represents a single chain configuration
type ChainConfig struct {
	Policy string       `yaml:"policy,omitempty"` // only for built-in chains
	Rules  []RuleConfig `yaml:"rules"`
}

// RuleConfig represents a single iptables rule
type RuleConfig struct {
	Protocol     string `yaml:"protocol,omitempty"`
	Source       string `yaml:"source,omitempty"`
	Destination  string `yaml:"destination,omitempty"`
	InInterface  string `yaml:"in-interface,omitempty"`
	OutInterface string `yaml:"out-interface,omitempty"`
	DPort        string `yaml:"dport,omitempty"`
	SPort        string `yaml:"sport,omitempty"`
	Jump         string `yaml:"jump"`

	Match []MatchConfig `yaml:"match,omitempty"`
}

// MatchConfig represents an iptables match module configuration
type MatchConfig struct {
	Name    string            `yaml:"name"`
	Options map[string]string `yaml:"options,omitempty"`
}
