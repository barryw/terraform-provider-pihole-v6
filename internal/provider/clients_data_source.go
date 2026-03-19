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
	_ datasource.DataSource              = &ClientsDataSource{}
	_ datasource.DataSourceWithConfigure = &ClientsDataSource{}
)

type ClientsDataSource struct {
	apiClient *pihole.Client
}

type ClientsDataSourceModel struct {
	Clients []ClientModel `tfsdk:"clients"`
}

type ClientModel struct {
	ID      types.String `tfsdk:"id"`
	Client  types.String `tfsdk:"client"`
	Comment types.String `tfsdk:"comment"`
	Groups  types.List   `tfsdk:"groups"`
}

func NewClientsDataSource() datasource.DataSource {
	return &ClientsDataSource{}
}

func (d *ClientsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_clients"
}

func (d *ClientsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches all clients.",
		Attributes: map[string]schema.Attribute{
			"clients": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":      schema.StringAttribute{Computed: true},
						"client":  schema.StringAttribute{Computed: true},
						"comment": schema.StringAttribute{Computed: true},
						"groups": schema.ListAttribute{
							Computed:    true,
							ElementType: types.Int64Type,
						},
					},
				},
			},
		},
	}
}

func (d *ClientsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ClientsDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	clients, err := d.apiClient.ListClients()
	if err != nil {
		resp.Diagnostics.AddError("Unable to list clients", err.Error())
		return
	}

	var state ClientsDataSourceModel
	for _, c := range clients {
		groupValues := make([]types.Int64, len(c.Groups))
		for i, g := range c.Groups {
			groupValues[i] = types.Int64Value(int64(g))
		}
		groupsList, diags := types.ListValueFrom(ctx, types.Int64Type, groupValues)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		state.Clients = append(state.Clients, ClientModel{
			ID:      types.StringValue(c.Client),
			Client:  types.StringValue(c.Client),
			Comment: types.StringValue(c.Comment),
			Groups:  groupsList,
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
