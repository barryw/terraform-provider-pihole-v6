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
	_ resource.Resource                = &DNSUpstreamsResource{}
	_ resource.ResourceWithImportState = &DNSUpstreamsResource{}
	_ resource.ResourceWithConfigure   = &DNSUpstreamsResource{}
)

// DNSUpstreamsResource manages the ordered list of upstream DNS servers PiHole
// forwards queries to. It is a singleton: there is one upstream list per
// instance, so define this resource at most once.
type DNSUpstreamsResource struct {
	client *pihole.Client
}

type DNSUpstreamsResourceModel struct {
	ID        types.String `tfsdk:"id"`
	Upstreams types.List   `tfsdk:"upstreams"`
}

func NewDNSUpstreamsResource() resource.Resource {
	return &DNSUpstreamsResource{}
}

func (r *DNSUpstreamsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dns_upstreams"
}

func (r *DNSUpstreamsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages PiHole's upstream DNS servers (config.dns.upstreams). This is a singleton resource — define it at most once. Order is significant and preserved.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Always \"upstreams\".",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"upstreams": schema.ListAttribute{
				Description: "Ordered list of upstream DNS servers (e.g. \"1.1.1.1\", \"8.8.8.8#5353\", \"127.0.0.1#5335\").",
				Required:    true,
				ElementType: types.StringType,
			},
		},
	}
}

func (r *DNSUpstreamsResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *DNSUpstreamsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan DNSUpstreamsResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var upstreams []string
	resp.Diagnostics.Append(plan.Upstreams.ElementsAs(ctx, &upstreams, false)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.SetDNSUpstreams(ctx, upstreams); err != nil {
		resp.Diagnostics.AddError("Error setting DNS upstreams", err.Error())
		return
	}

	plan.ID = types.StringValue("upstreams")
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *DNSUpstreamsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state DNSUpstreamsResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	upstreams, err := r.client.GetDNSUpstreams(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error reading DNS upstreams", err.Error())
		return
	}
	list, diags := types.ListValueFrom(ctx, types.StringType, upstreams)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.Upstreams = list
	state.ID = types.StringValue("upstreams")
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *DNSUpstreamsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan DNSUpstreamsResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var upstreams []string
	resp.Diagnostics.Append(plan.Upstreams.ElementsAs(ctx, &upstreams, false)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.SetDNSUpstreams(ctx, upstreams); err != nil {
		resp.Diagnostics.AddError("Error setting DNS upstreams", err.Error())
		return
	}

	plan.ID = types.StringValue("upstreams")
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *DNSUpstreamsResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
	// PiHole always has at least an implicit upstream configuration; there is
	// nothing to delete. Removing the resource from state leaves the current
	// upstreams in place.
}

func (r *DNSUpstreamsResource) ImportState(ctx context.Context, _ resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), "upstreams")...)
}
