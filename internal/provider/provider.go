package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Ensure LazyProvider satisfies various provider interfaces.
var _ provider.Provider = &LazyProvider{}

type LazyProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and run locally, and "test" when running acceptance
	// testing.
	version string
}

type LazyProviderModel struct {}

func (p *LazyProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "lazy"
	resp.Version = p.version
}

func (p *LazyProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{}
}

func (p *LazyProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {}

func (p *LazyProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
    NewLazyStringResource,
  }
}

func (p *LazyProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return nil
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &LazyProvider{
			version: version,
		}
	}
}

