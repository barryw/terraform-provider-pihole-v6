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
	_ datasource.DataSource              = &AdlistDataSource{}
	_ datasource.DataSourceWithConfigure = &AdlistDataSource{}
)

type AdlistDataSource struct {
	client *pihole.Client
}

type AdlistDataSourceModel struct {
	ID      types.String `tfsdk:"id"`
	Address types.String `tfsdk:"address"`
	Type    types.String `tfsdk:"type"`
	Comment types.String `tfsdk:"comment"`
	Groups  types.List   `tfsdk:"groups"`
	Enabled types.Bool   `tfsdk:"enabled"`
}

func NewAdlistDataSource() datasource.DataSource {
	return &AdlistDataSource{}
}

func (d *AdlistDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_adlist"
}

func (d *AdlistDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a single adlist by address.",
		Attributes: map[string]schema.Attribute{
			"id":      schema.StringAttribute{Computed: true, Description: "The address of the adlist."},
			"address": schema.StringAttribute{Required: true, Description: "The URL of the adlist to look up."},
			"type":    schema.StringAttribute{Computed: true, Description: "The type of the adlist: 'block' or 'allow'."},
			"comment": schema.StringAttribute{Computed: true, Description: "The comment for the adlist."},
			"groups": schema.ListAttribute{
				Computed:    true,
				ElementType: types.Int64Type,
				Description: "List of group IDs this adlist is assigned to.",
			},
			"enabled": schema.BoolAttribute{Computed: true, Description: "Whether the adlist is enabled."},
		},
	}
}

func (d *AdlistDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*pihole.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *pihole.Client, got: %T", req.ProviderData),
		)
		return
	}
	d.client = client
}

func (d *AdlistDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config AdlistDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	adlist, err := d.client.GetAdlist(config.Address.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to read adlist", err.Error())
		return
	}

	config.ID = types.StringValue(adlist.Address)
	config.Address = types.StringValue(adlist.Address)
	config.Type = types.StringValue(adlist.Type)
	config.Comment = types.StringValue(adlist.Comment)
	config.Enabled = types.BoolValue(adlist.Enabled)
	config.Groups = intSliceToInt64List(adlist.Groups)

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
