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
	_ datasource.DataSource              = &DNSBlockingDataSource{}
	_ datasource.DataSourceWithConfigure = &DNSBlockingDataSource{}
)

type DNSBlockingDataSource struct {
	client *pihole.Client
}

type DNSBlockingDataSourceModel struct {
	ID      types.String  `tfsdk:"id"`
	Enabled types.Bool    `tfsdk:"enabled"`
	Timer   types.Float64 `tfsdk:"timer"`
}

func NewDNSBlockingDataSource() datasource.DataSource {
	return &DNSBlockingDataSource{}
}

func (d *DNSBlockingDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dns_blocking"
}

func (d *DNSBlockingDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Reads PiHole's current global DNS blocking state.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true},
			"enabled": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether DNS blocking is currently enabled.",
			},
			"timer": schema.Float64Attribute{
				Computed:    true,
				Description: "Seconds remaining until blocking reverts to its previous state, or null if no timer is active.",
			},
		},
	}
}

func (d *DNSBlockingDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DNSBlockingDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config DNSBlockingDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	status, err := d.client.GetDNSBlocking(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Unable to read DNS blocking state", err.Error())
		return
	}

	config.ID = types.StringValue("blocking")
	config.Enabled = types.BoolValue(status.Enabled)
	if status.Timer != nil {
		config.Timer = types.Float64Value(*status.Timer)
	} else {
		config.Timer = types.Float64Null()
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
