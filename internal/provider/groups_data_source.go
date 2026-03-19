package provider

import (
	"context"

	pihole "github.com/barryw/go-pihole"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &GroupsDataSource{}
	_ datasource.DataSourceWithConfigure = &GroupsDataSource{}
)

type GroupsDataSource struct {
	client *pihole.Client
}

type GroupsDataSourceModel struct {
	Groups []GroupModel `tfsdk:"groups"`
}

type GroupModel struct {
	ID      types.String `tfsdk:"id"`
	Name    types.String `tfsdk:"name"`
	Comment types.String `tfsdk:"comment"`
	Enabled types.Bool   `tfsdk:"enabled"`
}

func NewGroupsDataSource() datasource.DataSource { return &GroupsDataSource{} }

func (d *GroupsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_groups"
}

func (d *GroupsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists all PiHole groups.",
		Attributes: map[string]schema.Attribute{
			"groups": schema.ListNestedAttribute{
				Description: "List of all groups.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "The group name (same as name).",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "The name of the group.",
							Computed:    true,
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
				},
			},
		},
	}
}

func (d *GroupsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *GroupsDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	groups, err := d.client.ListGroups()
	if err != nil {
		resp.Diagnostics.AddError("Error listing groups", err.Error())
		return
	}

	var state GroupsDataSourceModel
	for _, g := range groups {
		state.Groups = append(state.Groups, GroupModel{
			ID:      types.StringValue(g.Name),
			Name:    types.StringValue(g.Name),
			Comment: types.StringValue(g.Comment),
			Enabled: types.BoolValue(g.Enabled),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
