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
	_ datasource.DataSource              = &DNSUpstreamsDataSource{}
	_ datasource.DataSourceWithConfigure = &DNSUpstreamsDataSource{}
)

type DNSUpstreamsDataSource struct {
	client *pihole.Client
}

type DNSUpstreamsDataSourceModel struct {
	ID        types.String `tfsdk:"id"`
	Upstreams types.List   `tfsdk:"upstreams"`
}

func NewDNSUpstreamsDataSource() datasource.DataSource {
	return &DNSUpstreamsDataSource{}
}

func (d *DNSUpstreamsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dns_upstreams"
}

func (d *DNSUpstreamsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Reads PiHole's current ordered list of upstream DNS servers.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Computed: true},
			"upstreams": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "Ordered list of configured upstream DNS servers.",
			},
		},
	}
}

func (d *DNSUpstreamsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DNSUpstreamsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config DNSUpstreamsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	upstreams, err := d.client.GetDNSUpstreams(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Unable to read DNS upstreams", err.Error())
		return
	}

	list, diags := types.ListValueFrom(ctx, types.StringType, upstreams)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	config.ID = types.StringValue("upstreams")
	config.Upstreams = list
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
