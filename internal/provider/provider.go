package provider

import (
	"context"
	"fmt"
	"os"

	pihole "github.com/barryw/go-pihole"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ provider.Provider = &PiholeProvider{}

type PiholeProvider struct {
	version string
}

type PiholeProviderModel struct {
	URL      types.String `tfsdk:"url"`
	Password types.String `tfsdk:"password"`
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &PiholeProvider{version: version}
	}
}

func (p *PiholeProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "pihole"
	resp.Version = p.version
}

func (p *PiholeProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage PiHole v6 configuration via its API.",
		Attributes: map[string]schema.Attribute{
			"url": schema.StringAttribute{
				Description: "Base URL of the PiHole instance (e.g. http://192.168.1.1:8080). May also be set via PIHOLE_URL env var.",
				Optional:    true,
			},
			"password": schema.StringAttribute{
				Description: "App-password for the PiHole API. May also be set via PIHOLE_PASSWORD env var.",
				Optional:    true,
				Sensitive:   true,
			},
		},
	}
}

func (p *PiholeProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config PiholeProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.URL.IsUnknown() {
		resp.Diagnostics.AddAttributeError(path.Root("url"), "Unknown PiHole URL",
			"The provider cannot create the API client as there is an unknown configuration value for the URL.")
	}
	if config.Password.IsUnknown() {
		resp.Diagnostics.AddAttributeError(path.Root("password"), "Unknown PiHole Password",
			"The provider cannot create the API client as there is an unknown configuration value for the password.")
	}
	if resp.Diagnostics.HasError() {
		return
	}

	url := os.Getenv("PIHOLE_URL")
	password := os.Getenv("PIHOLE_PASSWORD")

	if !config.URL.IsNull() {
		url = config.URL.ValueString()
	}
	if !config.Password.IsNull() {
		password = config.Password.ValueString()
	}

	if url == "" {
		resp.Diagnostics.AddAttributeError(path.Root("url"), "Missing PiHole URL",
			"Set url in provider config or PIHOLE_URL environment variable.")
		return
	}
	if password == "" {
		resp.Diagnostics.AddAttributeError(path.Root("password"), "Missing PiHole Password",
			"Set password in provider config or PIHOLE_PASSWORD environment variable.")
		return
	}

	client, err := pihole.NewClient(url, password)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create PiHole client", err.Error())
		return
	}

	// Authenticate eagerly to avoid rate limiting when Terraform
	// makes many parallel requests with the same provider instance.
	if err := client.Authenticate(); err != nil {
		resp.Diagnostics.AddError("Failed to authenticate with PiHole",
			fmt.Sprintf("URL: %s — %s", url, err.Error()))
		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *PiholeProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewDNSRecordResource,
		NewCNAMERecordResource,
		NewGroupResource,
		NewAdlistResource,
		NewDomainListResource,
		NewClientResource,
		NewSettingResource,
	}
}

func (p *PiholeProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewDNSRecordDataSource,
		NewDNSRecordsDataSource,
		NewCNAMERecordDataSource,
		NewCNAMERecordsDataSource,
		NewGroupDataSource,
		NewGroupsDataSource,
		NewAdlistDataSource,
		NewAdlistsDataSource,
		NewDomainListDataSource,
		NewDomainListsDataSource,
		NewClientDataSource,
		NewClientsDataSource,
		NewSettingDataSource,
	}
}
