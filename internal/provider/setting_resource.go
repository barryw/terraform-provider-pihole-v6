package provider

import (
	"context"
	"encoding/json"
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
	_ resource.Resource                = &SettingResource{}
	_ resource.ResourceWithImportState = &SettingResource{}
	_ resource.ResourceWithConfigure   = &SettingResource{}
)

type SettingResource struct {
	client *pihole.Client
}

type SettingResourceModel struct {
	ID    types.String `tfsdk:"id"`
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

func NewSettingResource() resource.Resource {
	return &SettingResource{}
}

func (r *SettingResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_setting"
}

func (r *SettingResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a PiHole configuration setting. Values are JSON-encoded — use jsonencode() in your Terraform config.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Same as key.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"key": schema.StringAttribute{
				Description: "Dot-notation config path (e.g. webserver.api.app_sudo).",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"value": schema.StringAttribute{
				Description: "JSON-encoded value. Use jsonencode(true), jsonencode(86400), jsonencode(\"NULL\"), etc.",
				Required:    true,
			},
		},
	}
}

func (r *SettingResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *SettingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan SettingResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.SetConfig(plan.Key.ValueString(), json.RawMessage(plan.Value.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError("Error setting config", err.Error())
		return
	}

	plan.ID = plan.Key
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SettingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state SettingResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	val, err := r.client.GetConfig(state.Key.ValueString())
	if err != nil {
		if _, ok := err.(*pihole.ErrNotFound); ok {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading config", err.Error())
		return
	}

	state.Value = types.StringValue(string(val))
	state.ID = state.Key
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *SettingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan SettingResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.SetConfig(plan.Key.ValueString(), json.RawMessage(plan.Value.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError("Error setting config", err.Error())
		return
	}

	plan.ID = plan.Key
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SettingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// PiHole settings can't be truly deleted — they always have a value.
	// On destroy, we just remove from state. The setting remains at its current value.
	// Users who want to reset to default should set the default value before destroying.
}

func (r *SettingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("key"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}
