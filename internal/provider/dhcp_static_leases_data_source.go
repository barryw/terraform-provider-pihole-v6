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
	_ datasource.DataSource              = &DHCPStaticLeasesDataSource{}
	_ datasource.DataSourceWithConfigure = &DHCPStaticLeasesDataSource{}
)

type DHCPStaticLeasesDataSource struct {
	client *pihole.Client
}

type DHCPStaticLeasesDataSourceModel struct {
	Leases []DHCPStaticLeaseModel `tfsdk:"leases"`
}

type DHCPStaticLeaseModel struct {
	ID       types.String `tfsdk:"id"`
	MAC      types.String `tfsdk:"mac"`
	IP       types.String `tfsdk:"ip"`
	Hostname types.String `tfsdk:"hostname"`
}

func NewDHCPStaticLeasesDataSource() datasource.DataSource {
	return &DHCPStaticLeasesDataSource{}
}

func (d *DHCPStaticLeasesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dhcp_static_leases"
}

func (d *DHCPStaticLeasesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches all DHCP static leases.",
		Attributes: map[string]schema.Attribute{
			"leases": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":       schema.StringAttribute{Computed: true},
						"mac":      schema.StringAttribute{Computed: true, Description: "The MAC address of the device."},
						"ip":       schema.StringAttribute{Computed: true, Description: "The IP address assigned to the device."},
						"hostname": schema.StringAttribute{Computed: true, Description: "The hostname assigned to the device."},
					},
				},
			},
		},
	}
}

func (d *DHCPStaticLeasesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DHCPStaticLeasesDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	leases, err := d.client.ListDHCPStaticLeases()
	if err != nil {
		resp.Diagnostics.AddError("Unable to list DHCP static leases", err.Error())
		return
	}

	var state DHCPStaticLeasesDataSourceModel
	for _, l := range leases {
		state.Leases = append(state.Leases, DHCPStaticLeaseModel{
			ID:       types.StringValue(fmt.Sprintf("%s:%s", l.MAC, l.IP)),
			MAC:      types.StringValue(l.MAC),
			IP:       types.StringValue(l.IP),
			Hostname: types.StringValue(l.Hostname),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
