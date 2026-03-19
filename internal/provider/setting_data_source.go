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
	_ datasource.DataSource              = &SettingDataSource{}
	_ datasource.DataSourceWithConfigure = &SettingDataSource{}
)

type SettingDataSource struct {
	client *pihole.Client
}

type SettingDataSourceModel struct {
	ID    types.String `tfsdk:"id"`
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

func NewSettingDataSource() datasource.DataSource {
	return &SettingDataSource{}
}

func (d *SettingDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_setting"
}

func (d *SettingDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Reads a PiHole configuration setting. The value is returned as a JSON-encoded string — use jsondecode() to parse it.",
		Attributes: map[string]schema.Attribute{
			"id":  schema.StringAttribute{Computed: true},
			"key": schema.StringAttribute{Required: true, Description: "Dot-notation config path (e.g. webserver.api.app_sudo)."},
			"value": schema.StringAttribute{
				Computed:    true,
				Description: "JSON-encoded current value. Use jsondecode() to parse.",
			},
		},
	}
}

func (d *SettingDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *SettingDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config SettingDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	val, err := d.client.GetConfig(config.Key.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to read config setting", err.Error())
		return
	}

	config.ID = config.Key
	config.Value = types.StringValue(string(val))
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
