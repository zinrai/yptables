package generator

import (
	"fmt"
	"strings"

	"github.com/zinrai/yptables/internal/config"
)

// Generator generates iptables-restore format output
type Generator struct{}

// New creates a new Generator
func New() *Generator {
	return &Generator{}
}

// Generate generates iptables-restore format from the configuration
func (g *Generator) Generate(cfg *config.Config) ([]string, error) {
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

func (g *Generator) generateRuleLine(chain string, rule config.RuleConfig) string {
	parts := []string{"-A", chain}

	if rule.Protocol != "" {
		parts = append(parts, "-p", rule.Protocol)
	}
	if rule.Source != "" {
		parts = append(parts, "-s", rule.Source)
	}
	if rule.Destination != "" {
		parts = append(parts, "-d", rule.Destination)
	}
	if rule.InInterface != "" {
		parts = append(parts, "-i", rule.InInterface)
	}
	if rule.OutInterface != "" {
		parts = append(parts, "-o", rule.OutInterface)
	}
	if rule.DPort != "" {
		parts = append(parts, "--dport", rule.DPort)
	}
	if rule.SPort != "" {
		parts = append(parts, "--sport", rule.SPort)
	}

	for _, match := range rule.Match {
		parts = append(parts, "-m", match.Name)
		for optName, optValue := range match.Options {
			parts = append(parts, fmt.Sprintf("--%s", optName), formatOptionValue(optName, optValue))
		}
	}

	if rule.Jump != "" {
		parts = append(parts, "-j", rule.Jump)
	}

	return strings.Join(parts, " ")
}

func isBuiltinChain(tableName, chainName string) bool {
	builtinChains := map[string][]string{
		"filter": {"INPUT", "FORWARD", "OUTPUT"},
		"nat":    {"PREROUTING", "POSTROUTING", "OUTPUT"},
	}
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

func formatOptionValue(name, value string) string {
	if name == "comment" && strings.Contains(value, " ") {
		return fmt.Sprintf("\"%s\"", value)
	}
	return value
}
