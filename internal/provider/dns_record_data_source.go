package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

type DNSRecordDataSource struct{}

func NewDNSRecordDataSource() datasource.DataSource { return &DNSRecordDataSource{} }

func (d *DNSRecordDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dns_record"
}
func (d *DNSRecordDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{}
}
func (d *DNSRecordDataSource) Read(_ context.Context, _ datasource.ReadRequest, _ *datasource.ReadResponse) {
}
