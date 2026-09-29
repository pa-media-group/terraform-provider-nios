package dns

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// In every update plan dnssec_keys is unknown (computed-only, no
// UseStateForUnknown). It must not appear in the request: "dnssec_keys": []
// asks NIOS to delete the zone's DNSSEC keys ("Cannot delete the active key.").
func TestZoneAuthExpandOmitsUnknownDnssecKeys(t *testing.T) {
	var diags diag.Diagnostics
	m := ZoneAuthModel{
		Fqdn:       types.StringValue("example1.com"),
		View:       types.StringValue("default"),
		DnssecKeys: types.ListUnknown(types.ObjectType{AttrTypes: ZoneAuthDnssecKeysAttrTypes}),
	}
	payload := m.Expand(context.Background(), &diags, false)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var fields map[string]any
	if err := json.Unmarshal(body, &fields); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if v, ok := fields["dnssec_keys"]; ok {
		t.Fatalf("update payload contains dnssec_keys = %v; it must be omitted when unknown", v)
	}
}
