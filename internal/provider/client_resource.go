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
	_ resource.Resource                = &ClientResource{}
	_ resource.ResourceWithImportState = &ClientResource{}
	_ resource.ResourceWithConfigure   = &ClientResource{}
)

type ClientResource struct {
	apiClient *pihole.Client
}

type ClientResourceModel struct {
	ID      types.String `tfsdk:"id"`
	Client  types.String `tfsdk:"client"`
	Comment types.String `tfsdk:"comment"`
	Groups  types.List   `tfsdk:"groups"`
}

func NewClientResource() resource.Resource {
	return &ClientResource{}
}

func (r *ClientResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_client"
}

func (r *ClientResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a client in PiHole.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The client identifier.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"client": schema.StringAttribute{
				Description: "The client identifier (IP, MAC, CIDR, or hostname).",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"comment": schema.StringAttribute{
				Description: "Optional comment for the client.",
				Optional:    true,
				Computed:    true,
			},
			"groups": schema.ListAttribute{
				Description: "List of group IDs the client belongs to.",
				Optional:    true,
				Computed:    true,
				ElementType: types.Int64Type,
			},
		},
	}
}

func (r *ClientResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*pihole.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *pihole.Client, got: %T", req.ProviderData))
		return
	}
	r.apiClient = client
}

func (r *ClientResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ClientResourceModel
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

	created, err := r.apiClient.CreateClient(pihole.ClientCreateRequest{
		Client:  plan.Client.ValueString(),
		Comment: plan.Comment.ValueString(),
		Groups:  groups,
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating client", err.Error())
		return
	}

	plan.ID = types.StringValue(created.Client)
	plan.Client = types.StringValue(created.Client)
	plan.Comment = types.StringValue(created.Comment)

	groupValues := make([]types.Int64, len(created.Groups))
	for i, g := range created.Groups {
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

func (r *ClientResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ClientResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.apiClient.GetClient(state.Client.ValueString())
	if err != nil {
		if _, ok := err.(*pihole.ErrNotFound); ok {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading client", err.Error())
		return
	}

	state.ID = types.StringValue(client.Client)
	state.Client = types.StringValue(client.Client)
	state.Comment = types.StringValue(client.Comment)

	groupValues := make([]types.Int64, len(client.Groups))
	for i, g := range client.Groups {
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

func (r *ClientResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ClientResourceModel
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

	updated, err := r.apiClient.UpdateClient(plan.Client.ValueString(), pihole.ClientUpdateRequest{
		Comment: plan.Comment.ValueString(),
		Groups:  groups,
	})
	if err != nil {
		resp.Diagnostics.AddError("Error updating client", err.Error())
		return
	}

	plan.ID = types.StringValue(updated.Client)
	plan.Client = types.StringValue(updated.Client)
	plan.Comment = types.StringValue(updated.Comment)

	groupValues := make([]types.Int64, len(updated.Groups))
	for i, g := range updated.Groups {
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

func (r *ClientResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ClientResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.apiClient.DeleteClient(state.Client.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting client", err.Error())
		return
	}
}

func (r *ClientResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("client"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}
