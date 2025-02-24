package generator

import (
	"fmt"
	"strings"

	"github.com/zinrai/yptables/internal/config"
)

// Format specifies the output format
type Format int

const (
	ShellScript Format = iota
	IPTablesRestore
)

// Generator generates iptables commands
type Generator struct {
	format Format
}

// New creates a new Generator
func New(format Format) *Generator {
	return &Generator{format: format}
}

// Generate generates iptables commands from the configuration
func (g *Generator) Generate(config *config.Config) ([]string, error) {
	if g.format == ShellScript {
		return g.generateShellScript(config)
	}
	return g.generateIPTablesRestore(config)
}

func (g *Generator) generateShellScript(config *config.Config) ([]string, error) {
	var commands []string
	commands = append(commands, "#!/bin/sh")
	commands = append(commands, "")

	for tableName, table := range config.Tables {
		// Set default policies for built-in chains
		for chainName, chain := range table.Chains {
			if chain.Policy != "" {
				cmd := fmt.Sprintf("iptables -t %s -P %s %s", tableName, chainName, chain.Policy)
				commands = append(commands, cmd)
			}
		}

		// Create custom chains
		for chainName := range table.Chains {
			if !isBuiltinChain(tableName, chainName) {
				cmd := fmt.Sprintf("iptables -t %s -N %s", tableName, chainName)
				commands = append(commands, cmd)
			}
		}

		// Add rules
		for chainName, chain := range table.Chains {
			for _, rule := range chain.Rules {
				cmd := g.generateRuleCommand(tableName, chainName, rule)
				commands = append(commands, cmd)
			}
		}
	}

	return commands, nil
}

func (g *Generator) generateIPTablesRestore(config *config.Config) ([]string, error) {
	var lines []string

	for tableName, table := range config.Tables {
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
				line := g.generateRestoreRuleLine(chainName, rule)
				lines = append(lines, line)
			}
		}

		lines = append(lines, "COMMIT")
		lines = append(lines, "")
	}

	return lines, nil
}

func (g *Generator) generateRuleCommand(table, chain string, rule config.RuleConfig) string {
	parts := []string{"iptables", "-t", table, "-A", chain}

	// 基本オプションの追加
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

func (g *Generator) generateRestoreRuleLine(chain string, rule config.RuleConfig) string {
	parts := []string{"-A", chain}

	// 基本オプションの追加
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
