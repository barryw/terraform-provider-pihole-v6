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
	_ datasource.DataSource              = &ClientDataSource{}
	_ datasource.DataSourceWithConfigure = &ClientDataSource{}
)

type ClientDataSource struct {
	apiClient *pihole.Client
}

type ClientDataSourceModel struct {
	ID      types.String `tfsdk:"id"`
	Client  types.String `tfsdk:"client"`
	Comment types.String `tfsdk:"comment"`
	Groups  types.List   `tfsdk:"groups"`
}

func NewClientDataSource() datasource.DataSource {
	return &ClientDataSource{}
}

func (d *ClientDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_client"
}

func (d *ClientDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a single client by its identifier.",
		Attributes: map[string]schema.Attribute{
			"id":      schema.StringAttribute{Computed: true, Description: "The client identifier."},
			"client":  schema.StringAttribute{Required: true, Description: "The client identifier (IP, MAC, CIDR, or hostname) to look up."},
			"comment": schema.StringAttribute{Computed: true, Description: "Comment for the client."},
			"groups": schema.ListAttribute{
				Computed:    true,
				ElementType: types.Int64Type,
				Description: "List of group IDs the client belongs to.",
			},
		},
	}
}

func (d *ClientDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*pihole.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *pihole.Client, got: %T", req.ProviderData))
		return
	}
	d.apiClient = client
}

func (d *ClientDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config ClientDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := d.apiClient.GetClient(config.Client.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to read client", err.Error())
		return
	}

	config.ID = types.StringValue(client.Client)
	config.Client = types.StringValue(client.Client)
	config.Comment = types.StringValue(client.Comment)

	groupValues := make([]types.Int64, len(client.Groups))
	for i, g := range client.Groups {
		groupValues[i] = types.Int64Value(int64(g))
	}
	groupsList, diags := types.ListValueFrom(ctx, types.Int64Type, groupValues)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	config.Groups = groupsList

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
