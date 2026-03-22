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
	_ resource.Resource                = &DHCPStaticLeaseResource{}
	_ resource.ResourceWithImportState = &DHCPStaticLeaseResource{}
	_ resource.ResourceWithConfigure   = &DHCPStaticLeaseResource{}
)

type DHCPStaticLeaseResource struct {
	client *pihole.Client
}

type DHCPStaticLeaseResourceModel struct {
	ID       types.String `tfsdk:"id"`
	MAC      types.String `tfsdk:"mac"`
	IP       types.String `tfsdk:"ip"`
	Hostname types.String `tfsdk:"hostname"`
}

func NewDHCPStaticLeaseResource() resource.Resource {
	return &DHCPStaticLeaseResource{}
}

func (r *DHCPStaticLeaseResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dhcp_static_lease"
}

func (r *DHCPStaticLeaseResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a DHCP static lease (reservation) in PiHole.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Composite ID in format mac:ip.",
				Computed:    true,
			},
			"mac": schema.StringAttribute{
				Description: "The MAC address of the device.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"ip": schema.StringAttribute{
				Description: "The IP address to assign to the device.",
				Required:    true,
			},
			"hostname": schema.StringAttribute{
				Description: "Optional hostname for the device.",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func (r *DHCPStaticLeaseResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *DHCPStaticLeaseResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan DHCPStaticLeaseResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	lease := pihole.DHCPStaticLease{
		MAC:      plan.MAC.ValueString(),
		IP:       plan.IP.ValueString(),
		Hostname: plan.Hostname.ValueString(),
	}

	if err := r.client.CreateDHCPStaticLease(lease); err != nil {
		resp.Diagnostics.AddError("Error creating DHCP static lease", err.Error())
		return
	}

	plan.ID = types.StringValue(fmt.Sprintf("%s:%s", plan.MAC.ValueString(), plan.IP.ValueString()))
	if plan.Hostname.IsNull() || plan.Hostname.IsUnknown() {
		plan.Hostname = types.StringValue("")
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *DHCPStaticLeaseResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state DHCPStaticLeaseResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	lease, err := r.client.GetDHCPStaticLease(state.MAC.ValueString())
	if err != nil {
		if _, ok := err.(*pihole.ErrNotFound); ok {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading DHCP static lease", err.Error())
		return
	}

	state.MAC = types.StringValue(lease.MAC)
	state.IP = types.StringValue(lease.IP)
	state.Hostname = types.StringValue(lease.Hostname)
	state.ID = types.StringValue(fmt.Sprintf("%s:%s", lease.MAC, lease.IP))
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *DHCPStaticLeaseResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan DHCPStaticLeaseResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state DHCPStaticLeaseResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	oldLease := pihole.DHCPStaticLease{
		MAC:      state.MAC.ValueString(),
		IP:       state.IP.ValueString(),
		Hostname: state.Hostname.ValueString(),
	}
	newLease := pihole.DHCPStaticLease{
		MAC:      plan.MAC.ValueString(),
		IP:       plan.IP.ValueString(),
		Hostname: plan.Hostname.ValueString(),
	}

	if err := r.client.UpdateDHCPStaticLease(oldLease, newLease); err != nil {
		resp.Diagnostics.AddError("Error updating DHCP static lease", err.Error())
		return
	}

	plan.ID = types.StringValue(fmt.Sprintf("%s:%s", plan.MAC.ValueString(), plan.IP.ValueString()))
	if plan.Hostname.IsNull() || plan.Hostname.IsUnknown() {
		plan.Hostname = types.StringValue("")
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *DHCPStaticLeaseResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state DHCPStaticLeaseResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	lease := pihole.DHCPStaticLease{
		MAC:      state.MAC.ValueString(),
		IP:       state.IP.ValueString(),
		Hostname: state.Hostname.ValueString(),
	}

	if err := r.client.DeleteDHCPStaticLease(lease); err != nil {
		resp.Diagnostics.AddError("Error deleting DHCP static lease", err.Error())
		return
	}
}

func (r *DHCPStaticLeaseResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.SplitN(req.ID, ":", 7) // MAC has 5 colons, then :IP
	if len(parts) < 7 {
		resp.Diagnostics.AddError("Invalid import ID",
			fmt.Sprintf("Expected format: mac:ip (e.g. aa:bb:cc:dd:ee:ff:192.168.1.1). Got: %q", req.ID))
		return
	}
	mac := strings.Join(parts[:6], ":")
	ip := parts[6]
	if mac == "" || ip == "" {
		resp.Diagnostics.AddError("Invalid import ID",
			fmt.Sprintf("Expected format: mac:ip (e.g. aa:bb:cc:dd:ee:ff:192.168.1.1). Got: %q", req.ID))
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("mac"), mac)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("ip"), ip)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), fmt.Sprintf("%s:%s", mac, ip))...)
}
