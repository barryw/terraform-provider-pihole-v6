package provider

import (
	"context"
	"fmt"
	"strings"

	pihole "github.com/barryw/go-pihole"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &DNSRecordResource{}
	_ resource.ResourceWithImportState = &DNSRecordResource{}
	_ resource.ResourceWithConfigure   = &DNSRecordResource{}
)

type DNSRecordResource struct {
	client *pihole.Client
}

type DNSRecordResourceModel struct {
	ID     types.String `tfsdk:"id"`
	Domain types.String `tfsdk:"domain"`
	IP     types.String `tfsdk:"ip"`
}

func NewDNSRecordResource() resource.Resource {
	return &DNSRecordResource{}
}

func (r *DNSRecordResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dns_record"
}

func (r *DNSRecordResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a local DNS record in PiHole. Multiple records for the same domain with different IPs are supported.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Composite ID in format domain:ip.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"domain": schema.StringAttribute{
				Description: "The domain name.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"ip": schema.StringAttribute{
				Description: "The IP address to resolve to.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *DNSRecordResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *DNSRecordResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan DNSRecordResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.CreateDNSRecord(plan.IP.ValueString(), plan.Domain.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error creating DNS record", err.Error())
		return
	}

	plan.ID = types.StringValue(fmt.Sprintf("%s:%s", plan.Domain.ValueString(), plan.IP.ValueString()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *DNSRecordResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state DNSRecordResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// List all records and match by both domain AND IP to support multiple
	// records for the same domain with different IPs (round-robin DNS).
	records, err := r.client.ListDNSRecords()
	if err != nil {
		resp.Diagnostics.AddError("Error reading DNS records", err.Error())
		return
	}

	var found bool
	for _, r := range records {
		if r.Domain == state.Domain.ValueString() && r.IP == state.IP.ValueString() {
			found = true
			break
		}
	}

	if !found {
		resp.State.RemoveResource(ctx)
		return
	}

	state.ID = types.StringValue(fmt.Sprintf("%s:%s", state.Domain.ValueString(), state.IP.ValueString()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *DNSRecordResource) Update(_ context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Update not supported", "DNS records are immutable. Changes require replacement.")
}

func (r *DNSRecordResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state DNSRecordResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteDNSRecord(state.IP.ValueString(), state.Domain.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting DNS record", err.Error())
		return
	}
}

func (r *DNSRecordResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.SplitN(req.ID, ":", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError("Invalid import ID",
			fmt.Sprintf("Expected format: domain:ip. Got: %q", req.ID))
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("domain"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("ip"), parts[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}
