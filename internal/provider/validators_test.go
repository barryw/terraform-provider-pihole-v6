package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestStringOneOf(t *testing.T) {
	v := stringOneOf("allow", "deny")

	tests := []struct {
		name    string
		value   types.String
		wantErr bool
	}{
		{"valid allow", types.StringValue("allow"), false},
		{"valid deny", types.StringValue("deny"), false},
		{"invalid", types.StringValue("nope"), true},
		{"null skipped", types.StringNull(), false},
		{"unknown skipped", types.StringUnknown(), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validator.StringRequest{ConfigValue: tt.value}
			resp := &validator.StringResponse{}
			v.ValidateString(context.Background(), req, resp)
			if resp.Diagnostics.HasError() != tt.wantErr {
				t.Errorf("value %v: hasError=%v, want %v (%v)", tt.value, resp.Diagnostics.HasError(), tt.wantErr, resp.Diagnostics)
			}
		})
	}
}
