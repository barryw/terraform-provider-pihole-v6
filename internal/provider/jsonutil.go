package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// canonicalizeJSON validates that s is a single well-formed JSON value and
// returns its canonical encoding (compact, sorted object keys). Storing the
// canonical form keeps Terraform state consistent with what the PiHole API
// returns on Read, avoiding spurious perpetual diffs when users write
// non-canonical JSON.
func canonicalizeJSON(s string) (string, error) {
	dec := json.NewDecoder(bytes.NewReader([]byte(s)))
	// UseNumber keeps JSON numbers as their original literal instead of
	// coercing to float64, which would lose precision for integers > 2^53.
	dec.UseNumber()
	var v interface{}
	if err := dec.Decode(&v); err != nil {
		return "", fmt.Errorf("value must be valid JSON: %w", err)
	}
	// Reject trailing content after the first JSON value (e.g. "true false").
	if dec.More() {
		return "", fmt.Errorf("value must be a single JSON value")
	}
	out, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(out), nil
}
