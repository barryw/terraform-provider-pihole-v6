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
	_ datasource.DataSource              = &AdlistsDataSource{}
	_ datasource.DataSourceWithConfigure = &AdlistsDataSource{}
)

type AdlistsDataSource struct {
	client *pihole.Client
}

type AdlistsDataSourceModel struct {
	Adlists []AdlistEntryModel `tfsdk:"adlists"`
}

type AdlistEntryModel struct {
	ID      types.String `tfsdk:"id"`
	Address types.String `tfsdk:"address"`
	Type    types.String `tfsdk:"type"`
	Comment types.String `tfsdk:"comment"`
	Groups  types.List   `tfsdk:"groups"`
	Enabled types.Bool   `tfsdk:"enabled"`
}

func NewAdlistsDataSource() datasource.DataSource {
	return &AdlistsDataSource{}
}

func (d *AdlistsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_adlists"
}

func (d *AdlistsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches all adlists from PiHole.",
		Attributes: map[string]schema.Attribute{
			"adlists": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of all adlists.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":      schema.StringAttribute{Computed: true, Description: "The address of the adlist."},
						"address": schema.StringAttribute{Computed: true, Description: "The URL of the adlist."},
						"type":    schema.StringAttribute{Computed: true, Description: "The type of the adlist: 'block' or 'allow'."},
						"comment": schema.StringAttribute{Computed: true, Description: "The comment for the adlist."},
						"groups": schema.ListAttribute{
							Computed:    true,
							ElementType: types.Int64Type,
							Description: "List of group IDs this adlist is assigned to.",
						},
						"enabled": schema.BoolAttribute{Computed: true, Description: "Whether the adlist is enabled."},
					},
				},
			},
		},
	}
}

func (d *AdlistsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *AdlistsDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	adlists, err := d.client.ListAdlists()
	if err != nil {
		resp.Diagnostics.AddError("Unable to list adlists", err.Error())
		return
	}

	var state AdlistsDataSourceModel
	state.Adlists = make([]AdlistEntryModel, len(adlists))
	for i, a := range adlists {
		state.Adlists[i] = AdlistEntryModel{
			ID:      types.StringValue(a.Address),
			Address: types.StringValue(a.Address),
			Type:    types.StringValue(a.Type),
			Comment: types.StringValue(a.Comment),
			Enabled: types.BoolValue(a.Enabled),
			Groups:  intSliceToInt64List(a.Groups),
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
