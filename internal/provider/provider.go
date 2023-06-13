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

type LazyProviderModel struct{}

func (p *LazyProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "lazy"
	resp.Version = p.version
}

func (p *LazyProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: `
Maintaining a Terraform configuration so that it can encode several different deployment contexts is difficult.  Some contexts may be aware of values that others are not: If you run Terraform in CI/CD pipelines and one pipeline does some preliminary calculation to set a variable, other pipelines or attempts to run the Terraform locally may be difficult to implement, just because they need to work out a meaningful value to set.

The ` + "`lazy`" + ` provider is an attempt at a solution to this issue; it provides a resource that will track a value and only update its result if a new value is explicitly given. In the hypothetical situation mentioned above, one pipeline might always set a value for a specific variable, while some may always leave it null. instead of Terraform requiring the value to be set, the previous value will automatically be used, and there will be no changes to resources that depend on that variable.
    `,
	}
}

func (p *LazyProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
}

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
