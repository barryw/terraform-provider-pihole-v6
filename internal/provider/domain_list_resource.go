package provider

import (
	"context"
	"fmt"
	"strings"

	pihole "github.com/barryw/go-pihole"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &DomainListResource{}
	_ resource.ResourceWithImportState = &DomainListResource{}
	_ resource.ResourceWithConfigure   = &DomainListResource{}
)

type DomainListResource struct {
	client *pihole.Client
}

type DomainListResourceModel struct {
	ID      types.String `tfsdk:"id"`
	Domain  types.String `tfsdk:"domain"`
	Type    types.String `tfsdk:"type"`
	Kind    types.String `tfsdk:"kind"`
	Comment types.String `tfsdk:"comment"`
	Groups  types.List   `tfsdk:"groups"`
	Enabled types.Bool   `tfsdk:"enabled"`
}

func NewDomainListResource() resource.Resource {
	return &DomainListResource{}
}

func (r *DomainListResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_domain_list"
}

func (r *DomainListResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a domain list entry (allow/deny, exact/regex) in PiHole.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Composite ID in format type:kind:domain.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"domain": schema.StringAttribute{
				Description: "The domain or regex pattern.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				Description: "The list type: allow or deny.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"kind": schema.StringAttribute{
				Description: "The match kind: exact or regex.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"comment": schema.StringAttribute{
				Description: "Optional comment for this domain entry.",
				Optional:    true,
				Computed:    true,
			},
			"groups": schema.ListAttribute{
				Description: "List of group IDs this domain entry belongs to.",
				Optional:    true,
				Computed:    true,
				ElementType: types.Int64Type,
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether this domain entry is enabled.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
		},
	}
}

func (r *DomainListResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *DomainListResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan DomainListResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var groups []int
	if !plan.Groups.IsNull() && !plan.Groups.IsUnknown() {
		var g []int64
		resp.Diagnostics.Append(plan.Groups.ElementsAs(ctx, &g, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		for _, v := range g {
			groups = append(groups, int(v))
		}
	}

	createReq := pihole.DomainCreateRequest{
		Domain:  plan.Domain.ValueString(),
		Type:    plan.Type.ValueString(),
		Kind:    plan.Kind.ValueString(),
		Comment: plan.Comment.ValueString(),
		Groups:  groups,
		Enabled: plan.Enabled.ValueBool(),
	}

	entry, err := r.client.CreateDomain(createReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating domain list entry", err.Error())
		return
	}

	plan.ID = types.StringValue(fmt.Sprintf("%s:%s:%s", entry.Type, entry.Kind, entry.Domain))
	plan.Domain = types.StringValue(entry.Domain)
	plan.Type = types.StringValue(entry.Type)
	plan.Kind = types.StringValue(entry.Kind)
	plan.Comment = types.StringValue(entry.Comment)
	plan.Enabled = types.BoolValue(entry.Enabled)

	groupValues := make([]types.Int64, len(entry.Groups))
	for i, g := range entry.Groups {
		groupValues[i] = types.Int64Value(int64(g))
	}
	groupsList, diags := types.ListValueFrom(ctx, types.Int64Type, groupValues)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.Groups = groupsList

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *DomainListResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state DomainListResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	entry, err := r.client.GetDomain(state.Type.ValueString(), state.Kind.ValueString(), state.Domain.ValueString())
	if err != nil {
		if _, ok := err.(*pihole.ErrNotFound); ok {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading domain list entry", err.Error())
		return
	}

	state.ID = types.StringValue(fmt.Sprintf("%s:%s:%s", entry.Type, entry.Kind, entry.Domain))
	state.Domain = types.StringValue(entry.Domain)
	state.Type = types.StringValue(entry.Type)
	state.Kind = types.StringValue(entry.Kind)
	state.Comment = types.StringValue(entry.Comment)
	state.Enabled = types.BoolValue(entry.Enabled)

	groupValues := make([]types.Int64, len(entry.Groups))
	for i, g := range entry.Groups {
		groupValues[i] = types.Int64Value(int64(g))
	}
	groupsList, diags := types.ListValueFrom(ctx, types.Int64Type, groupValues)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.Groups = groupsList

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *DomainListResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan DomainListResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var groups []int
	if !plan.Groups.IsNull() && !plan.Groups.IsUnknown() {
		var g []int64
		resp.Diagnostics.Append(plan.Groups.ElementsAs(ctx, &g, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		for _, v := range g {
			groups = append(groups, int(v))
		}
	}

	updateReq := pihole.DomainUpdateRequest{
		Type:    plan.Type.ValueString(),
		Kind:    plan.Kind.ValueString(),
		Comment: plan.Comment.ValueString(),
		Groups:  groups,
		Enabled: plan.Enabled.ValueBool(),
	}

	entry, err := r.client.UpdateDomain(plan.Type.ValueString(), plan.Kind.ValueString(), plan.Domain.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Error updating domain list entry", err.Error())
		return
	}

	plan.ID = types.StringValue(fmt.Sprintf("%s:%s:%s", entry.Type, entry.Kind, entry.Domain))
	plan.Domain = types.StringValue(entry.Domain)
	plan.Type = types.StringValue(entry.Type)
	plan.Kind = types.StringValue(entry.Kind)
	plan.Comment = types.StringValue(entry.Comment)
	plan.Enabled = types.BoolValue(entry.Enabled)

	groupValues := make([]types.Int64, len(entry.Groups))
	for i, g := range entry.Groups {
		groupValues[i] = types.Int64Value(int64(g))
	}
	groupsList, diags := types.ListValueFrom(ctx, types.Int64Type, groupValues)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.Groups = groupsList

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *DomainListResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state DomainListResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteDomain(state.Type.ValueString(), state.Kind.ValueString(), state.Domain.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting domain list entry", err.Error())
		return
	}
}

func (r *DomainListResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.SplitN(req.ID, ":", 3)
	if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
		resp.Diagnostics.AddError("Invalid import ID",
			fmt.Sprintf("Expected format: type:kind:domain. Got: %q", req.ID))
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("type"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("kind"), parts[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("domain"), parts[2])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}
