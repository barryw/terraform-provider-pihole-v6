package provider

import (
	"context"
	"fmt"

	pihole "github.com/barryw/go-pihole"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &DNSBlockingResource{}
	_ resource.ResourceWithImportState = &DNSBlockingResource{}
	_ resource.ResourceWithConfigure   = &DNSBlockingResource{}
)

// DNSBlockingResource manages PiHole's global DNS blocking on/off state. It is a
// singleton: there is exactly one blocking state per PiHole instance, so only
// one instance of this resource should exist. The temporary auto-revert timer
// is intentionally not modelled here because a countdown is imperative state,
// not declarative configuration.
type DNSBlockingResource struct {
	client *pihole.Client
}

type DNSBlockingResourceModel struct {
	ID      types.String `tfsdk:"id"`
	Enabled types.Bool   `tfsdk:"enabled"`
}

func NewDNSBlockingResource() resource.Resource {
	return &DNSBlockingResource{}
}

func (r *DNSBlockingResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dns_blocking"
}

func (r *DNSBlockingResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages PiHole's global DNS blocking state (enabled/disabled). This is a singleton resource — define it at most once. Destroying it stops managing the state but leaves blocking as-is.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Always \"blocking\".",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether DNS blocking is enabled.",
				Required:    true,
			},
		},
	}
}

func (r *DNSBlockingResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*pihole.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *pihole.Client, got: %T", req.ProviderData))
		return
	}
	r.client = client
}

func (r *DNSBlockingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan DNSBlockingResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.SetDNSBlocking(ctx, plan.Enabled.ValueBool(), nil); err != nil {
		resp.Diagnostics.AddError("Error setting DNS blocking", err.Error())
		return
	}
	plan.ID = types.StringValue("blocking")
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *DNSBlockingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state DNSBlockingResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	status, err := r.client.GetDNSBlocking(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error reading DNS blocking", err.Error())
		return
	}
	state.Enabled = types.BoolValue(status.Enabled)
	state.ID = types.StringValue("blocking")
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *DNSBlockingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan DNSBlockingResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.SetDNSBlocking(ctx, plan.Enabled.ValueBool(), nil); err != nil {
		resp.Diagnostics.AddError("Error setting DNS blocking", err.Error())
		return
	}
	plan.ID = types.StringValue("blocking")
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *DNSBlockingResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
	// Blocking is a global toggle that always has a value; there is nothing to
	// delete. Removing the resource from state leaves blocking in its current
	// state.
}

func (r *DNSBlockingResource) ImportState(ctx context.Context, _ resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), "blocking")...)
}
