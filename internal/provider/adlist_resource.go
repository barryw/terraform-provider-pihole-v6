package provider

import (
	"context"
	"errors"
	"fmt"

	pihole "github.com/barryw/go-pihole"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &AdlistResource{}
	_ resource.ResourceWithConfigure   = &AdlistResource{}
	_ resource.ResourceWithImportState = &AdlistResource{}
)

type AdlistResource struct {
	client *pihole.Client
}

type AdlistResourceModel struct {
	ID      types.String `tfsdk:"id"`
	Address types.String `tfsdk:"address"`
	Type    types.String `tfsdk:"type"`
	Comment types.String `tfsdk:"comment"`
	Groups  types.List   `tfsdk:"groups"`
	Enabled types.Bool   `tfsdk:"enabled"`
}

func NewAdlistResource() resource.Resource {
	return &AdlistResource{}
}

func (r *AdlistResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_adlist"
}

func (r *AdlistResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an adlist (block or allow list) in PiHole.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The address of the adlist (same as address).",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"address": schema.StringAttribute{
				Description: "The URL of the adlist.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				Description: "The type of the adlist: 'block' or 'allow'.",
				Required:    true,
			},
			"comment": schema.StringAttribute{
				Description: "An optional comment for the adlist.",
				Optional:    true,
				Computed:    true,
			},
			"groups": schema.ListAttribute{
				Description: "List of group IDs this adlist is assigned to.",
				Optional:    true,
				Computed:    true,
				ElementType: types.Int64Type,
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether the adlist is enabled.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
		},
	}
}

func (r *AdlistResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *AdlistResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan AdlistResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := pihole.AdlistCreateRequest{
		Address: plan.Address.ValueString(),
		Type:    plan.Type.ValueString(),
		Comment: plan.Comment.ValueString(),
		Groups:  int64ListToIntSlice(plan.Groups),
		Enabled: plan.Enabled.ValueBool(),
	}

	adlist, err := r.client.CreateAdlist(createReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating adlist", err.Error())
		return
	}

	mapAdlistToModel(adlist, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *AdlistResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state AdlistResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	adlist, err := r.client.GetAdlist(state.Address.ValueString())
	if err != nil {
		var notFound *pihole.ErrNotFound
		if errors.As(err, &notFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading adlist", err.Error())
		return
	}

	mapAdlistToModel(adlist, &state)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *AdlistResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan AdlistResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateReq := pihole.AdlistUpdateRequest{
		Comment: plan.Comment.ValueString(),
		Type:    plan.Type.ValueString(),
		Groups:  int64ListToIntSlice(plan.Groups),
		Enabled: plan.Enabled.ValueBool(),
	}

	adlist, err := r.client.UpdateAdlist(plan.Address.ValueString(), plan.Type.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Error updating adlist", err.Error())
		return
	}

	mapAdlistToModel(adlist, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *AdlistResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state AdlistResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteAdlist(state.Address.ValueString(), state.Type.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting adlist", err.Error())
		return
	}
}

func (r *AdlistResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("address"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}

// mapAdlistToModel maps an API Adlist response to the Terraform resource model.
func mapAdlistToModel(adlist *pihole.Adlist, model *AdlistResourceModel) {
	model.ID = types.StringValue(adlist.Address)
	model.Address = types.StringValue(adlist.Address)
	model.Type = types.StringValue(adlist.Type)
	model.Comment = types.StringValue(adlist.Comment)
	model.Enabled = types.BoolValue(adlist.Enabled)
	model.Groups = intSliceToInt64List(adlist.Groups)
}

// intSliceToInt64List converts []int to a types.List of Int64 values.
func intSliceToInt64List(ints []int) types.List {
	elems := make([]attr.Value, len(ints))
	for i, v := range ints {
		elems[i] = types.Int64Value(int64(v))
	}
	return types.ListValueMust(types.Int64Type, elems)
}

// int64ListToIntSlice converts a types.List of Int64 values to []int.
func int64ListToIntSlice(list types.List) []int {
	if list.IsNull() || list.IsUnknown() {
		return nil
	}
	elems := list.Elements()
	result := make([]int, len(elems))
	for i, elem := range elems {
		result[i] = int(elem.(types.Int64).ValueInt64())
	}
	return result
}
