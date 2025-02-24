package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

// LoadFromFile reads and parses the YAML configuration file
func LoadFromFile(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("parsing config file: %w", err)
	}

	if err := validate(&config); err != nil {
		return nil, fmt.Errorf("validating config: %w", err)
	}

	return &config, nil
}

// validate checks if the configuration is valid
func validate(config *Config) error {
	for tableName, table := range config.Tables {
		if err := validateTable(tableName, table); err != nil {
			return err
		}
	}

	return nil
}

// validateTable checks if the table configuration is valid
func validateTable(name string, table *TableConfig) error {
	if name != "filter" && name != "nat" {
		return fmt.Errorf("unsupported table: %s", name)
	}

	for chainName, chain := range table.Chains {
		if err := validateChain(name, chainName, chain); err != nil {
			return err
		}
	}

	return nil
}

// validateChain checks if the chain configuration is valid
func validateChain(tableName, chainName string, chain *ChainConfig) error {
	// Check if built-in chain has policy
	builtinChains := getBuiltinChains(tableName)
	isBuiltin := false
	for _, c := range builtinChains {
		if chainName == c {
			isBuiltin = true
			break
		}
	}

	if isBuiltin && chain.Policy != "" {
		if chain.Policy != "ACCEPT" && chain.Policy != "DROP" {
			return fmt.Errorf("invalid policy for chain %s: %s", chainName, chain.Policy)
		}
	}

	return nil
}

// getBuiltinChains returns the built-in chains for a given table
func getBuiltinChains(tableName string) []string {
	switch tableName {
	case "filter":
		return []string{"INPUT", "FORWARD", "OUTPUT"}
	case "nat":
		return []string{"PREROUTING", "POSTROUTING", "OUTPUT"}
	default:
		return []string{}
	}
}
