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
	_ datasource.DataSource              = &CNAMERecordDataSource{}
	_ datasource.DataSourceWithConfigure = &CNAMERecordDataSource{}
)

type CNAMERecordDataSource struct {
	client *pihole.Client
}

type CNAMERecordDataSourceModel struct {
	ID     types.String `tfsdk:"id"`
	Domain types.String `tfsdk:"domain"`
	Target types.String `tfsdk:"target"`
	TTL    types.Int64  `tfsdk:"ttl"`
}

func NewCNAMERecordDataSource() datasource.DataSource { return &CNAMERecordDataSource{} }

func (d *CNAMERecordDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cname_record"
}

func (d *CNAMERecordDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a single CNAME record by domain.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Composite ID in the format domain:target.",
				Computed:    true,
			},
			"domain": schema.StringAttribute{
				Description: "The domain name to look up.",
				Required:    true,
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
	}
}

func (d *CNAMERecordDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *CNAMERecordDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config CNAMERecordDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	record, err := d.client.GetCNAMERecord(config.Domain.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading CNAME record", err.Error())
		return
	}

	config.Domain = types.StringValue(record.Domain)
	config.Target = types.StringValue(record.Target)
	config.TTL = types.Int64Value(int64(record.TTL))
	config.ID = types.StringValue(record.Domain + ":" + record.Target)

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
