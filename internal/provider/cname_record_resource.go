package provider

import (
	"context"
	"errors"
	"fmt"
	"strings"

	pihole "github.com/barryw/go-pihole"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &CNAMERecordResource{}
	_ resource.ResourceWithConfigure   = &CNAMERecordResource{}
	_ resource.ResourceWithImportState = &CNAMERecordResource{}
)

type CNAMERecordResource struct {
	client *pihole.Client
}

type CNAMERecordResourceModel struct {
	ID     types.String `tfsdk:"id"`
	Domain types.String `tfsdk:"domain"`
	Target types.String `tfsdk:"target"`
	TTL    types.Int64  `tfsdk:"ttl"`
}

func NewCNAMERecordResource() resource.Resource { return &CNAMERecordResource{} }

func (r *CNAMERecordResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cname_record"
}

func (r *CNAMERecordResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a CNAME record in PiHole.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Composite ID in the format domain:target.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"domain": schema.StringAttribute{
				Description: "The domain name for the CNAME record.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"target": schema.StringAttribute{
				Description: "The target hostname for the CNAME record.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"ttl": schema.Int64Attribute{
				Description: "Time to live in seconds. 0 means use default.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *CNAMERecordResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*pihole.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *pihole.Client, got %T.", req.ProviderData),
		)
		return
	}
	r.client = client
}

func (r *CNAMERecordResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan CNAMERecordResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	domain := plan.Domain.ValueString()
	target := plan.Target.ValueString()
	ttl := int(plan.TTL.ValueInt64())

	if err := r.client.CreateCNAMERecord(domain, target, ttl); err != nil {
		resp.Diagnostics.AddError("Error creating CNAME record", err.Error())
		return
	}

	plan.ID = types.StringValue(domain + ":" + target)
	if plan.TTL.IsNull() || plan.TTL.IsUnknown() {
		plan.TTL = types.Int64Value(0)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *CNAMERecordResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state CNAMERecordResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	record, err := r.client.GetCNAMERecord(state.Domain.ValueString())
	if err != nil {
		var notFound *pihole.ErrNotFound
		if errors.As(err, &notFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading CNAME record", err.Error())
		return
	}

	state.Domain = types.StringValue(record.Domain)
	state.Target = types.StringValue(record.Target)
	state.TTL = types.Int64Value(int64(record.TTL))
	state.ID = types.StringValue(record.Domain + ":" + record.Target)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *CNAMERecordResource) Update(_ context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"CNAME records are immutable",
		"All attributes use RequiresReplace; Update should never be called.",
	)
}

func (r *CNAMERecordResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state CNAMERecordResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	domain := state.Domain.ValueString()
	target := state.Target.ValueString()
	ttl := int(state.TTL.ValueInt64())

	if err := r.client.DeleteCNAMERecord(domain, target, ttl); err != nil {
		resp.Diagnostics.AddError("Error deleting CNAME record", err.Error())
		return
	}
}

func (r *CNAMERecordResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.SplitN(req.ID, ":", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError(
			"Invalid import ID",
			"Expected format: domain:target",
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("domain"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("target"), parts[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}
