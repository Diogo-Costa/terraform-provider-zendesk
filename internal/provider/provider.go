package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ provider.Provider = &ZendeskProvider{}

type ZendeskProvider struct {
	version string
}

type ZendeskProviderModel struct {
	Subdomain types.String `tfsdk:"subdomain"`
	Email     types.String `tfsdk:"email"`
	APIToken  types.String `tfsdk:"api_token"`
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &ZendeskProvider{
			version: version,
		}
	}
}

func (p *ZendeskProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "zendesk"
	resp.Version = p.version
}

func (p *ZendeskProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Interact with Zendesk.",
		Attributes: map[string]schema.Attribute{
			"subdomain": schema.StringAttribute{
				Description: "The Zendesk subdomain (e.g., company in company.zendesk.com)",
				Required:    true,
			},
			"email": schema.StringAttribute{
				Description: "The email address associated with the Zendesk account",
				Required:    true,
			},
			"api_token": schema.StringAttribute{
				Description: "The API token for authentication",
				Required:    true,
				Sensitive:   true,
			},
		},
	}
}

func (p *ZendeskProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config ZendeskProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Subdomain.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("subdomain"),
			"Unknown Zendesk subdomain",
			"The provider cannot create the Zendesk API client as the subdomain is unknown.",
		)
	}

	if config.Email.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("email"),
			"Unknown Zendesk email",
			"The provider cannot create the Zendesk API client as the email is unknown.",
		)
	}

	if config.APIToken.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_token"),
			"Unknown Zendesk API token",
			"The provider cannot create the Zendesk API client as the API token is unknown.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	subdomain := os.Getenv("ZENDESK_SUBDOMAIN")
	email := os.Getenv("ZENDESK_EMAIL")
	apiToken := os.Getenv("ZENDESK_API_TOKEN")

	if !config.Subdomain.IsNull() {
		subdomain = config.Subdomain.ValueString()
	}

	if !config.Email.IsNull() {
		email = config.Email.ValueString()
	}

	if !config.APIToken.IsNull() {
		apiToken = config.APIToken.ValueString()
	}

	if subdomain == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("subdomain"),
			"Missing Zendesk subdomain",
			"The provider cannot create the Zendesk API client as the subdomain is missing.",
		)
	}

	if email == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("email"),
			"Missing Zendesk email",
			"The provider cannot create the Zendesk API client as the email is missing.",
		)
	}

	if apiToken == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_token"),
			"Missing Zendesk API token",
			"The provider cannot create the Zendesk API client as the API token is missing.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: Create Zendesk client
	// client := NewClient(subdomain, email, apiToken)
	// resp.DataSourceData = client
	// resp.ResourceData = client
}

func (p *ZendeskProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		// Add data sources here
	}
}

func (p *ZendeskProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewOAuthClientResource,
		NewOAuthTokenResource,
	}
} 