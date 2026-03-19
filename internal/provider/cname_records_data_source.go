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
	_ datasource.DataSource              = &CNAMERecordsDataSource{}
	_ datasource.DataSourceWithConfigure = &CNAMERecordsDataSource{}
)

type CNAMERecordsDataSource struct {
	client *pihole.Client
}

type CNAMERecordsDataSourceModel struct {
	Records []CNAMERecordModel `tfsdk:"records"`
}

type CNAMERecordModel struct {
	ID     types.String `tfsdk:"id"`
	Domain types.String `tfsdk:"domain"`
	Target types.String `tfsdk:"target"`
	TTL    types.Int64  `tfsdk:"ttl"`
}

func NewCNAMERecordsDataSource() datasource.DataSource { return &CNAMERecordsDataSource{} }

func (d *CNAMERecordsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cname_records"
}

func (d *CNAMERecordsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches all CNAME records from PiHole.",
		Attributes: map[string]schema.Attribute{
			"records": schema.ListNestedAttribute{
				Description: "List of CNAME records.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Composite ID in the format domain:target.",
							Computed:    true,
						},
						"domain": schema.StringAttribute{
							Description: "The domain name for the CNAME record.",
							Computed:    true,
						},
						"target": schema.StringAttribute{
							Description: "The target hostname for the CNAME record.",
							Computed:    true,
						},
						"ttl": schema.Int64Attribute{
							Description: "Time to live in seconds.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *CNAMERecordsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*pihole.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected DataSource Configure Type",
			fmt.Sprintf("Expected *pihole.Client, got %T.", req.ProviderData),
		)
		return
	}
	d.client = client
}

func (d *CNAMERecordsDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	records, err := d.client.ListCNAMERecords()
	if err != nil {
		resp.Diagnostics.AddError("Error listing CNAME records", err.Error())
		return
	}

	var state CNAMERecordsDataSourceModel
	state.Records = make([]CNAMERecordModel, len(records))
	for i, r := range records {
		state.Records[i] = CNAMERecordModel{
			ID:     types.StringValue(r.Domain + ":" + r.Target),
			Domain: types.StringValue(r.Domain),
			Target: types.StringValue(r.Target),
			TTL:    types.Int64Value(int64(r.TTL)),
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
