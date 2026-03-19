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
	_ datasource.DataSource              = &DomainListsDataSource{}
	_ datasource.DataSourceWithConfigure = &DomainListsDataSource{}
)

type DomainListsDataSource struct {
	client *pihole.Client
}

type DomainListsDataSourceModel struct {
	Type    types.String            `tfsdk:"type"`
	Kind    types.String            `tfsdk:"kind"`
	Domains []DomainListEntryModel  `tfsdk:"domains"`
}

type DomainListEntryModel struct {
	ID      types.String `tfsdk:"id"`
	Domain  types.String `tfsdk:"domain"`
	Type    types.String `tfsdk:"type"`
	Kind    types.String `tfsdk:"kind"`
	Comment types.String `tfsdk:"comment"`
	Groups  types.List   `tfsdk:"groups"`
	Enabled types.Bool   `tfsdk:"enabled"`
}

func NewDomainListsDataSource() datasource.DataSource {
	return &DomainListsDataSource{}
}

func (d *DomainListsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_domain_lists"
}

func (d *DomainListsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches all domain list entries, with optional type and kind filters.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Optional:    true,
				Description: "Filter by list type: allow or deny.",
			},
			"kind": schema.StringAttribute{
				Optional:    true,
				Description: "Filter by match kind: exact or regex.",
			},
			"domains": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of domain list entries.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":      schema.StringAttribute{Computed: true, Description: "Composite ID in format type:kind:domain."},
						"domain":  schema.StringAttribute{Computed: true, Description: "The domain or regex pattern."},
						"type":    schema.StringAttribute{Computed: true, Description: "The list type: allow or deny."},
						"kind":    schema.StringAttribute{Computed: true, Description: "The match kind: exact or regex."},
						"comment": schema.StringAttribute{Computed: true, Description: "Comment for this domain entry."},
						"groups": schema.ListAttribute{
							Computed:    true,
							ElementType: types.Int64Type,
							Description: "List of group IDs this domain entry belongs to.",
						},
						"enabled": schema.BoolAttribute{Computed: true, Description: "Whether this domain entry is enabled."},
					},
				},
			},
		},
	}
}

func (d *DomainListsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DomainListsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config DomainListsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var entries []pihole.DomainEntry
	var err error

	if !config.Type.IsNull() && !config.Kind.IsNull() {
		entries, err = d.client.ListDomainsByTypeAndKind(config.Type.ValueString(), config.Kind.ValueString())
	} else {
		entries, err = d.client.ListDomains()
	}
	if err != nil {
		resp.Diagnostics.AddError("Unable to list domain entries", err.Error())
		return
	}

	// Apply client-side filtering if only one of type/kind is set
	if !config.Type.IsNull() && config.Kind.IsNull() {
		filtered := make([]pihole.DomainEntry, 0, len(entries))
		for _, e := range entries {
			if e.Type == config.Type.ValueString() {
				filtered = append(filtered, e)
			}
		}
		entries = filtered
	} else if config.Type.IsNull() && !config.Kind.IsNull() {
		filtered := make([]pihole.DomainEntry, 0, len(entries))
		for _, e := range entries {
			if e.Kind == config.Kind.ValueString() {
				filtered = append(filtered, e)
			}
		}
		entries = filtered
	}

	domains := make([]DomainListEntryModel, len(entries))
	for i, entry := range entries {
		groupValues := make([]types.Int64, len(entry.Groups))
		for j, g := range entry.Groups {
			groupValues[j] = types.Int64Value(int64(g))
		}
		groupsList, diags := types.ListValueFrom(ctx, types.Int64Type, groupValues)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		domains[i] = DomainListEntryModel{
			ID:      types.StringValue(fmt.Sprintf("%s:%s:%s", entry.Type, entry.Kind, entry.Domain)),
			Domain:  types.StringValue(entry.Domain),
			Type:    types.StringValue(entry.Type),
			Kind:    types.StringValue(entry.Kind),
			Comment: types.StringValue(entry.Comment),
			Groups:  groupsList,
			Enabled: types.BoolValue(entry.Enabled),
		}
	}

	config.Domains = domains
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
