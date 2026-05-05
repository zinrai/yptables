package main

import (
	"fmt"
	"maps"
	"slices"
)

// Generator generates iptables-restore format output
type Generator struct{}

// NewGenerator creates a new Generator
func NewGenerator() *Generator {
	return &Generator{}
}

// Generate generates iptables-restore format from the configuration.
// Tables, chains, and match options are emitted in deterministic alphabetical
// order so that the same input always produces byte-identical output.
func (g *Generator) Generate(cfg *Config) ([]string, error) {
	var lines []string

	for _, tableName := range slices.Sorted(maps.Keys(cfg.Tables)) {
		table := cfg.Tables[tableName]
		lines = append(lines, fmt.Sprintf("*%s", tableName))

		chainNames := slices.Sorted(maps.Keys(table.Chains))

		// Set default policies for built-in chains
		for _, chainName := range chainNames {
			chain := table.Chains[chainName]
			if chain.Policy != "" {
				lines = append(lines, fmt.Sprintf(":%s %s [0:0]", chainName, chain.Policy))
			}
		}

		// Create custom chains
		for _, chainName := range chainNames {
			if !isBuiltinChain(tableName, chainName) {
				lines = append(lines, fmt.Sprintf(":%s - [0:0]", chainName))
			}
		}

		// Add rules
		for _, chainName := range chainNames {
			chain := table.Chains[chainName]
			for _, rule := range chain.Rules {
				line := g.generateRuleLine(chainName, rule)
				lines = append(lines, line)
			}
		}

		lines = append(lines, "COMMIT")
		lines = append(lines, "")
	}

	return lines, nil
}
