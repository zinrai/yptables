package main

import (
	"strings"
	"testing"
)

func TestValidate_AcceptsSupportedTables(t *testing.T) {
	for _, table := range []string{"filter", "nat"} {
		t.Run(table, func(t *testing.T) {
			cfg := &Config{Tables: map[string]*TableConfig{
				table: {Chains: map[string]*ChainConfig{}},
			}}
			if err := validate(cfg); err != nil {
				t.Errorf("expected no error for %s, got: %v", table, err)
			}
		})
	}
}

func TestValidate_RejectsUnsupportedTable(t *testing.T) {
	cfg := &Config{Tables: map[string]*TableConfig{
		"mangle": {Chains: map[string]*ChainConfig{}},
	}}
	err := validate(cfg)
	if err == nil {
		t.Fatal("expected error for unsupported table 'mangle'")
	}
	if !strings.Contains(err.Error(), "mangle") {
		t.Errorf("error should mention 'mangle', got: %v", err)
	}
}

func TestValidate_BuiltinChainPolicy(t *testing.T) {
	tests := []struct {
		name    string
		table   string
		chain   string
		policy  string
		wantErr bool
	}{
		{"INPUT ACCEPT", "filter", "INPUT", "ACCEPT", false},
		{"INPUT DROP", "filter", "INPUT", "DROP", false},
		{"INPUT empty policy is allowed", "filter", "INPUT", "", false},
		{"INPUT REJECT is rejected", "filter", "INPUT", "REJECT", true},
		{"INPUT QUEUE is rejected", "filter", "INPUT", "QUEUE", true},
		{"PREROUTING ACCEPT", "nat", "PREROUTING", "ACCEPT", false},
		{"PREROUTING MASQUERADE is rejected", "nat", "PREROUTING", "MASQUERADE", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{Tables: map[string]*TableConfig{
				tt.table: {Chains: map[string]*ChainConfig{
					tt.chain: {Policy: tt.policy},
				}},
			}}
			err := validate(cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Pins current contract: policy on custom chains is silently accepted.
// Generate ignores it (only built-in chains emit a policy line).
// If we ever decide to reject it instead, this test should be flipped.
func TestValidate_CustomChainPolicyNotChecked(t *testing.T) {
	cfg := &Config{Tables: map[string]*TableConfig{
		"filter": {Chains: map[string]*ChainConfig{
			"my-custom-chain": {Policy: "REJECT"},
		}},
	}}
	if err := validate(cfg); err != nil {
		t.Errorf("custom chain policy should not be validated currently, got: %v", err)
	}
}
