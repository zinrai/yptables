package main

import (
	"slices"
	"strings"
	"testing"
)

func TestGenerate_BuiltinPolicyAndCustomChainDeclaration(t *testing.T) {
	cfg := &Config{Tables: map[string]*TableConfig{
		"filter": {Chains: map[string]*ChainConfig{
			"INPUT": {
				Policy: "DROP",
				Rules: []RuleConfig{
					{Protocol: "tcp", Jump: "WEB"},
				},
			},
			"WEB": {
				Rules: []RuleConfig{
					{Source: "10.0.0.0/8", Jump: "ACCEPT"},
				},
			},
		}},
	}}

	g := NewGenerator()
	got, err := g.Generate(cfg)
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}

	want := []string{
		"*filter",
		":INPUT DROP [0:0]",
		":WEB - [0:0]",
		"-A INPUT -p tcp -j WEB",
		"-A WEB -s 10.0.0.0/8 -j ACCEPT",
		"COMMIT",
		"",
	}
	if !slices.Equal(got, want) {
		t.Errorf("unexpected output:\ngot:\n%s\nwant:\n%s",
			strings.Join(got, "\n"), strings.Join(want, "\n"))
	}
}

func TestGenerate_TablesEmittedInAlphabeticalOrder(t *testing.T) {
	cfg := &Config{Tables: map[string]*TableConfig{
		"nat": {Chains: map[string]*ChainConfig{
			"PREROUTING": {Rules: []RuleConfig{{Protocol: "tcp", Jump: "ACCEPT"}}},
		}},
		"filter": {Chains: map[string]*ChainConfig{
			"INPUT": {Policy: "DROP"},
		}},
	}}

	g := NewGenerator()
	got, err := g.Generate(cfg)
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}

	want := []string{
		"*filter",
		":INPUT DROP [0:0]",
		"COMMIT",
		"",
		"*nat",
		"-A PREROUTING -p tcp -j ACCEPT",
		"COMMIT",
		"",
	}
	if !slices.Equal(got, want) {
		t.Errorf("unexpected output:\ngot:\n%s\nwant:\n%s",
			strings.Join(got, "\n"), strings.Join(want, "\n"))
	}
}

func TestGenerate_RulesPreserveSliceOrderWithinChain(t *testing.T) {
	cfg := &Config{Tables: map[string]*TableConfig{
		"filter": {Chains: map[string]*ChainConfig{
			"INPUT": {
				Policy: "DROP",
				Rules: []RuleConfig{
					{InInterface: "lo", Jump: "ACCEPT"},
					{Protocol: "icmp", Jump: "ACCEPT"},
					{Protocol: "tcp", DPort: "22", Jump: "ACCEPT"},
				},
			},
		}},
	}}

	g := NewGenerator()
	got, err := g.Generate(cfg)
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}

	want := []string{
		"*filter",
		":INPUT DROP [0:0]",
		"-A INPUT -i lo -j ACCEPT",
		"-A INPUT -p icmp -j ACCEPT",
		"-A INPUT -p tcp --dport 22 -j ACCEPT",
		"COMMIT",
		"",
	}
	if !slices.Equal(got, want) {
		t.Errorf("unexpected output:\ngot:\n%s\nwant:\n%s",
			strings.Join(got, "\n"), strings.Join(want, "\n"))
	}
}
