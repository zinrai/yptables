package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// builtinChains lists the built-in chains for each supported table.
var builtinChains = map[string][]string{
	"filter": {"INPUT", "FORWARD", "OUTPUT"},
	"nat":    {"PREROUTING", "POSTROUTING", "OUTPUT"},
}

// isBuiltinChain reports whether chainName is a built-in chain of tableName.
func isBuiltinChain(tableName, chainName string) bool {
	chains, ok := builtinChains[tableName]
	if !ok {
		return false
	}
	for _, c := range chains {
		if c == chainName {
			return true
		}
	}
	return false
}

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
	if isBuiltinChain(tableName, chainName) && chain.Policy != "" {
		if chain.Policy != "ACCEPT" && chain.Policy != "DROP" {
			return fmt.Errorf("invalid policy for chain %s: %s", chainName, chain.Policy)
		}
	}

	return nil
}
