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
	_ datasource.DataSource              = &DomainListDataSource{}
	_ datasource.DataSourceWithConfigure = &DomainListDataSource{}
)

type DomainListDataSource struct {
	client *pihole.Client
}

type DomainListDataSourceModel struct {
	ID      types.String `tfsdk:"id"`
	Domain  types.String `tfsdk:"domain"`
	Type    types.String `tfsdk:"type"`
	Kind    types.String `tfsdk:"kind"`
	Comment types.String `tfsdk:"comment"`
	Groups  types.List   `tfsdk:"groups"`
	Enabled types.Bool   `tfsdk:"enabled"`
}

func NewDomainListDataSource() datasource.DataSource {
	return &DomainListDataSource{}
}

func (d *DomainListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_domain_list"
}

func (d *DomainListDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches a single domain list entry by domain, type, and kind.",
		Attributes: map[string]schema.Attribute{
			"id":      schema.StringAttribute{Computed: true, Description: "Composite ID in format type:kind:domain."},
			"domain":  schema.StringAttribute{Required: true, Description: "The domain or regex pattern."},
			"type":    schema.StringAttribute{Required: true, Description: "The list type: allow or deny."},
			"kind":    schema.StringAttribute{Required: true, Description: "The match kind: exact or regex."},
			"comment": schema.StringAttribute{Computed: true, Description: "Comment for this domain entry."},
			"groups": schema.ListAttribute{
				Computed:    true,
				ElementType: types.Int64Type,
				Description: "List of group IDs this domain entry belongs to.",
			},
			"enabled": schema.BoolAttribute{Computed: true, Description: "Whether this domain entry is enabled."},
		},
	}
}

func (d *DomainListDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DomainListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config DomainListDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	entry, err := d.client.GetDomain(config.Type.ValueString(), config.Kind.ValueString(), config.Domain.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Unable to read domain list entry", err.Error())
		return
	}

	config.ID = types.StringValue(fmt.Sprintf("%s:%s:%s", entry.Type, entry.Kind, entry.Domain))
	config.Domain = types.StringValue(entry.Domain)
	config.Type = types.StringValue(entry.Type)
	config.Kind = types.StringValue(entry.Kind)
	config.Comment = types.StringValue(entry.Comment)
	config.Enabled = types.BoolValue(entry.Enabled)

	groupValues := make([]types.Int64, len(entry.Groups))
	for i, g := range entry.Groups {
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
