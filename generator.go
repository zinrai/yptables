package main

import "fmt"

// Generator generates iptables-restore format output
type Generator struct{}

// NewGenerator creates a new Generator
func NewGenerator() *Generator {
	return &Generator{}
}

// Generate generates iptables-restore format from the configuration
func (g *Generator) Generate(cfg *Config) ([]string, error) {
	var lines []string

	for tableName, table := range cfg.Tables {
		lines = append(lines, fmt.Sprintf("*%s", tableName))

		// Set default policies for built-in chains
		for chainName, chain := range table.Chains {
			if chain.Policy != "" {
				lines = append(lines, fmt.Sprintf(":%s %s [0:0]", chainName, chain.Policy))
			}
		}

		// Create custom chains
		for chainName := range table.Chains {
			if !isBuiltinChain(tableName, chainName) {
				lines = append(lines, fmt.Sprintf(":%s - [0:0]", chainName))
			}
		}

		// Add rules
		for chainName, chain := range table.Chains {
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
