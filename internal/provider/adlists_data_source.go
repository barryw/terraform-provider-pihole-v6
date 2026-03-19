package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

type AdlistsDataSource struct{}

func NewAdlistsDataSource() datasource.DataSource { return &AdlistsDataSource{} }

func (d *AdlistsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_adlists"
}
func (d *AdlistsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{}
}
func (d *AdlistsDataSource) Read(_ context.Context, _ datasource.ReadRequest, _ *datasource.ReadResponse) {
}
