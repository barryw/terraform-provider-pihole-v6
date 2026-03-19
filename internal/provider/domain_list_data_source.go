package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

type DomainListDataSource struct{}

func NewDomainListDataSource() datasource.DataSource { return &DomainListDataSource{} }

func (d *DomainListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_domain_list"
}
func (d *DomainListDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{}
}
func (d *DomainListDataSource) Read(_ context.Context, _ datasource.ReadRequest, _ *datasource.ReadResponse) {
}
