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
	_ datasource.DataSource              = &DNSRecordDataSource{}
	_ datasource.DataSourceWithConfigure = &DNSRecordDataSource{}
)

type DNSRecordDataSource struct {
	client *pihole.Client
}

type DNSRecordDataSourceModel struct {
	ID     types.String `tfsdk:"id"`
	Domain types.String `tfsdk:"domain"`
	IP     types.String `tfsdk:"ip"`
}

func NewDNSRecordDataSource() datasource.DataSource {
	return &DNSRecordDataSource{}
}

func (d *DNSRecordDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dns_record"
}

func (d *DNSRecordDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a single DNS record by domain name.",
		Attributes: map[string]schema.Attribute{
			"id":     schema.StringAttribute{Computed: true},
			"domain": schema.StringAttribute{Required: true, Description: "The domain name to look up."},
			"ip":     schema.StringAttribute{Computed: true, Description: "The IP address the domain resolves to."},
		},
	}
}

func (d *DNSRecordDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DNSRecordDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config DNSRecordDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	record, err := d.client.GetDNSRecord(config.Domain.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to read DNS record", err.Error())
		return
	}

	config.ID = types.StringValue(fmt.Sprintf("%s:%s", record.Domain, record.IP))
	config.IP = types.StringValue(record.IP)
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
