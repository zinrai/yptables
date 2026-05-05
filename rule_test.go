package main

import "testing"

func TestGenerateRuleLine(t *testing.T) {
	tests := []struct {
		name  string
		chain string
		rule  RuleConfig
		want  string
	}{
		{
			name:  "jump only",
			chain: "INPUT",
			rule:  RuleConfig{Jump: "ACCEPT"},
			want:  "-A INPUT -j ACCEPT",
		},
		{
			name:  "protocol mapping",
			chain: "INPUT",
			rule:  RuleConfig{Protocol: "tcp", Jump: "ACCEPT"},
			want:  "-A INPUT -p tcp -j ACCEPT",
		},
		{
			name:  "source mapping",
			chain: "INPUT",
			rule:  RuleConfig{Source: "10.0.0.0/8", Jump: "ACCEPT"},
			want:  "-A INPUT -s 10.0.0.0/8 -j ACCEPT",
		},
		{
			name:  "destination mapping",
			chain: "OUTPUT",
			rule:  RuleConfig{Destination: "192.168.0.1", Jump: "DROP"},
			want:  "-A OUTPUT -d 192.168.0.1 -j DROP",
		},
		{
			name:  "in-interface mapping",
			chain: "INPUT",
			rule:  RuleConfig{InInterface: "eth0", Jump: "ACCEPT"},
			want:  "-A INPUT -i eth0 -j ACCEPT",
		},
		{
			name:  "out-interface mapping",
			chain: "OUTPUT",
			rule:  RuleConfig{OutInterface: "eth1", Jump: "ACCEPT"},
			want:  "-A OUTPUT -o eth1 -j ACCEPT",
		},
		{
			name:  "dport mapping",
			chain: "INPUT",
			rule:  RuleConfig{DPort: "22", Jump: "ACCEPT"},
			want:  "-A INPUT --dport 22 -j ACCEPT",
		},
		{
			name:  "sport mapping",
			chain: "INPUT",
			rule:  RuleConfig{SPort: "1024", Jump: "ACCEPT"},
			want:  "-A INPUT --sport 1024 -j ACCEPT",
		},
		{
			name:  "all simple fields combined preserves canonical order",
			chain: "FORWARD",
			rule: RuleConfig{
				Protocol:     "tcp",
				Source:       "10.0.0.0/8",
				Destination:  "192.168.0.1",
				InInterface:  "eth0",
				OutInterface: "eth1",
				DPort:        "443",
				SPort:        "1024",
				Jump:         "ACCEPT",
			},
			want: "-A FORWARD -p tcp -s 10.0.0.0/8 -d 192.168.0.1 -i eth0 -o eth1 --dport 443 --sport 1024 -j ACCEPT",
		},
		{
			name:  "match with no options",
			chain: "INPUT",
			rule: RuleConfig{
				Match: []MatchConfig{{Name: "state"}},
				Jump:  "ACCEPT",
			},
			want: "-A INPUT -m state -j ACCEPT",
		},
		{
			name:  "match with single option",
			chain: "INPUT",
			rule: RuleConfig{
				Match: []MatchConfig{{Name: "tcp", Options: map[string]string{"dport": "22"}}},
				Jump:  "ACCEPT",
			},
			want: "-A INPUT -m tcp --dport 22 -j ACCEPT",
		},
		{
			name:  "multiple matches preserve YAML list order",
			chain: "INPUT",
			rule: RuleConfig{
				Protocol: "tcp",
				Match: []MatchConfig{
					{Name: "multiport", Options: map[string]string{"dports": "80,443"}},
					{Name: "comment", Options: map[string]string{"comment": "Web traffic"}},
				},
				Jump: "WEB",
			},
			want: `-A INPUT -p tcp -m multiport --dports 80,443 -m comment --comment "Web traffic" -j WEB`,
		},
		{
			name:  "comment with spaces is quoted",
			chain: "INPUT",
			rule: RuleConfig{
				Match: []MatchConfig{{Name: "comment", Options: map[string]string{"comment": "Allow internal"}}},
				Jump:  "ACCEPT",
			},
			want: `-A INPUT -m comment --comment "Allow internal" -j ACCEPT`,
		},
		{
			name:  "comment without spaces is not quoted",
			chain: "INPUT",
			rule: RuleConfig{
				Match: []MatchConfig{{Name: "comment", Options: map[string]string{"comment": "internal"}}},
				Jump:  "ACCEPT",
			},
			want: "-A INPUT -m comment --comment internal -j ACCEPT",
		},
		{
			name:  "no jump emits no -j",
			chain: "INPUT",
			rule:  RuleConfig{Protocol: "tcp"},
			want:  "-A INPUT -p tcp",
		},
		{
			name:  "match with multiple options emits them in alphabetical order",
			chain: "INPUT",
			rule: RuleConfig{
				Match: []MatchConfig{
					{Name: "conntrack", Options: map[string]string{
						"ctstate": "NEW",
						"ctdir":   "ORIGINAL",
					}},
				},
				Jump: "ACCEPT",
			},
			want: "-A INPUT -m conntrack --ctdir ORIGINAL --ctstate NEW -j ACCEPT",
		},
	}

	g := NewGenerator()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := g.generateRuleLine(tt.chain, tt.rule)
			if got != tt.want {
				t.Errorf("got:  %q\nwant: %q", got, tt.want)
			}
		})
	}
}

func TestFormatOptionValue(t *testing.T) {
	tests := []struct {
		name  string
		key   string
		value string
		want  string
	}{
		{"non-comment numeric", "dport", "22", "22"},
		{"non-comment with spaces is left as-is", "dports", "80, 443", "80, 443"},
		{"comment without spaces is unquoted", "comment", "internal", "internal"},
		{"comment with spaces is quoted", "comment", "Web traffic", `"Web traffic"`},
		{"empty value", "comment", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatOptionValue(tt.key, tt.value)
			if got != tt.want {
				t.Errorf("formatOptionValue(%q, %q) = %q, want %q",
					tt.key, tt.value, got, tt.want)
			}
		})
	}
}
