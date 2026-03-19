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
	_ datasource.DataSource              = &DNSRecordsDataSource{}
	_ datasource.DataSourceWithConfigure = &DNSRecordsDataSource{}
)

type DNSRecordsDataSource struct {
	client *pihole.Client
}

type DNSRecordsDataSourceModel struct {
	Records []DNSRecordModel `tfsdk:"records"`
}

type DNSRecordModel struct {
	ID     types.String `tfsdk:"id"`
	Domain types.String `tfsdk:"domain"`
	IP     types.String `tfsdk:"ip"`
}

func NewDNSRecordsDataSource() datasource.DataSource {
	return &DNSRecordsDataSource{}
}

func (d *DNSRecordsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dns_records"
}

func (d *DNSRecordsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches all DNS records.",
		Attributes: map[string]schema.Attribute{
			"records": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":     schema.StringAttribute{Computed: true},
						"domain": schema.StringAttribute{Computed: true},
						"ip":     schema.StringAttribute{Computed: true},
					},
				},
			},
		},
	}
}

func (d *DNSRecordsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DNSRecordsDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	records, err := d.client.ListDNSRecords()
	if err != nil {
		resp.Diagnostics.AddError("Unable to list DNS records", err.Error())
		return
	}

	var state DNSRecordsDataSourceModel
	for _, r := range records {
		state.Records = append(state.Records, DNSRecordModel{
			ID:     types.StringValue(fmt.Sprintf("%s:%s", r.Domain, r.IP)),
			Domain: types.StringValue(r.Domain),
			IP:     types.StringValue(r.IP),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
