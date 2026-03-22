package provider

import (
	"context"
	"fmt"

	pihole "github.com/barryw/go-pihole"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &DHCPStaticLeaseDataSource{}
	_ datasource.DataSourceWithConfigure = &DHCPStaticLeaseDataSource{}
)

type DHCPStaticLeaseDataSource struct {
	client *pihole.Client
}

type DHCPStaticLeaseDataSourceModel struct {
	ID       types.String `tfsdk:"id"`
	MAC      types.String `tfsdk:"mac"`
	IP       types.String `tfsdk:"ip"`
	Hostname types.String `tfsdk:"hostname"`
}

func NewDHCPStaticLeaseDataSource() datasource.DataSource {
	return &DHCPStaticLeaseDataSource{}
}

func (d *DHCPStaticLeaseDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dhcp_static_lease"
}

func (d *DHCPStaticLeaseDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a single DHCP static lease by MAC address.",
		Attributes: map[string]schema.Attribute{
			"id":       schema.StringAttribute{Computed: true},
			"mac":      schema.StringAttribute{Required: true, Description: "The MAC address to look up."},
			"ip":       schema.StringAttribute{Computed: true, Description: "The IP address assigned to the device."},
			"hostname": schema.StringAttribute{Computed: true, Description: "The hostname assigned to the device."},
		},
	}
}

func (d *DHCPStaticLeaseDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*pihole.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *pihole.Client, got: %T", req.ProviderData))
		return
	}
	d.client = client
}

func (d *DHCPStaticLeaseDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config DHCPStaticLeaseDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	lease, err := d.client.GetDHCPStaticLease(config.MAC.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to read DHCP static lease", err.Error())
		return
	}

	config.ID = types.StringValue(fmt.Sprintf("%s:%s", lease.MAC, lease.IP))
	config.IP = types.StringValue(lease.IP)
	config.Hostname = types.StringValue(lease.Hostname)
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
