package provider

import (
	"context"

	pihole "github.com/barryw/go-pihole"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &GroupDataSource{}
	_ datasource.DataSourceWithConfigure = &GroupDataSource{}
)

type GroupDataSource struct {
	client *pihole.Client
}

type GroupDataSourceModel struct {
	ID      types.String `tfsdk:"id"`
	Name    types.String `tfsdk:"name"`
	Comment types.String `tfsdk:"comment"`
	Enabled types.Bool   `tfsdk:"enabled"`
}

func NewGroupDataSource() datasource.DataSource { return &GroupDataSource{} }

func (d *GroupDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group"
}

func (d *GroupDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a single PiHole group by name.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The group name (same as name).",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the group to look up.",
				Required:    true,
			},
			"comment": schema.StringAttribute{
				Description: "The comment associated with the group.",
				Computed:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether the group is enabled.",
				Computed:    true,
			},
		},
	}
}

func (d *GroupDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*pihole.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Data Source Configure Type",
			"Expected *pihole.Client, got something else.")
		return
	}
	d.client = client
}

func (d *GroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config GroupDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	group, err := d.client.GetGroup(config.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading group", err.Error())
		return
	}

	config.ID = types.StringValue(group.Name)
	config.Name = types.StringValue(group.Name)
	config.Comment = types.StringValue(group.Comment)
	config.Enabled = types.BoolValue(group.Enabled)

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
