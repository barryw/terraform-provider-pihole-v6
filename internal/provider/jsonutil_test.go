package provider

import "testing"

func TestCanonicalizeJSON(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		want    string
		wantErr bool
	}{
		{"bool", "true", "true", false},
		{"int", "86400", "86400", false},
		{"string", `"NULL"`, `"NULL"`, false},
		{"leading space", "  true ", "true", false},
		{"object key order", `{"b":2,"a":1}`, `{"a":1,"b":2}`, false},
		{"large integer keeps precision", "1000000000000", "1000000000000", false},
		{"big int not float", "9007199254740993", "9007199254740993", false},
		{"invalid", "not json", "", true},
		{"empty", "", "", true},
		{"trailing garbage", "true false", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := canonicalizeJSON(tt.in)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error for %q", tt.in)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("canonicalizeJSON(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}
