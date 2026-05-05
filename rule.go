package main

import (
	"fmt"
	"maps"
	"slices"
	"strings"
)

func (g *Generator) generateRuleLine(chain string, rule RuleConfig) string {
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
		for _, optName := range slices.Sorted(maps.Keys(match.Options)) {
			optValue := match.Options[optName]
			parts = append(parts, fmt.Sprintf("--%s", optName), formatOptionValue(optName, optValue))
		}
	}

	if rule.Jump != "" {
		parts = append(parts, "-j", rule.Jump)
	}

	return strings.Join(parts, " ")
}

func formatOptionValue(name, value string) string {
	if name == "comment" && strings.Contains(value, " ") {
		return fmt.Sprintf("\"%s\"", value)
	}
	return value
}
