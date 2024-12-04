package provider

import (
	"context"
	"crypto/tls"
	"github.com/ccin2p3/go-freeipa/freeipa"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"net/http"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ provider.Provider = &freeipaProvider{}

var _ provider.ProviderWithFunctions = &freeipaProvider{}

type freeipaProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

type freeipaProviderModel struct {
	Host     types.String `tfsdk:"host"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
	Realm    types.String `tfsdk:"realm"`
	Insecure types.Bool   `tfsdk:"insecure"`
}

func (p *freeipaProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "freeipa"
	resp.Version = p.version
}

func (p *freeipaProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				MarkdownDescription: "The hostname of the FreeIPA master to use",
				Required:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "The username to use to authenticate with the FreeIPA master",
				Required:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "The password to use to authenticate with the FreeIPA master",
				Required:            true,
				Sensitive:           true,
			},
			"realm": schema.StringAttribute{
				MarkdownDescription: "The realm to use to authenticate with the FreeIPA master",
				Required:            true,
			},
			"insecure": schema.BoolAttribute{
				MarkdownDescription: "Whether to skip verification of the FreeIPA master's TLS certificate",
				Optional:            true,
			},
		},
	}
}

func (p *freeipaProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config freeipaProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"host is required",
			"host is required for the provider to function",
		)
	}

	if config.Username.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"username is required",
			"username is required for the provider to function",
		)
	}

	if config.Password.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"password is required",
			"password is required for the provider to function",
		)
	}

	if config.Realm.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("realm"),
			"realm is required",
			"realm is required for the provider to function",
		)
	}

	if config.Insecure.IsUnknown() {
		config.Insecure = types.BoolValue(false)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	host := os.Getenv("FREEIPA_HOST")
	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	username := os.Getenv("FREEIPA_USERNAME")
	if !config.Username.IsNull() {
		username = config.Username.ValueString()
	}

	password := os.Getenv("FREEIPA_PASSWORD")
	if !config.Password.IsNull() {
		password = config.Password.ValueString()
	}

	realm := os.Getenv("FREEIPA_REALM")
	if !config.Realm.IsNull() {
		realm = config.Realm.ValueString()
	}

	if host == "" {
		resp.Diagnostics.AddError("host is required",
			"host is required for the provider to function")
	}

	if username == "" {
		resp.Diagnostics.AddError("username is required",
			"username is required for the provider to function")
	}

	if password == "" {
		resp.Diagnostics.AddError("password is required",
			"password is required for the provider to function")
	}

	if realm == "" {
		resp.Diagnostics.AddError("realm is required",
			"realm is required for the provider to function")
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Create new freeipa client and set it as the data source and resource data
	tspt := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: config.Insecure.ValueBool(),
		},
	}
	client, err := freeipa.Connect(host, tspt, username, password)
	if err != nil {
		resp.Diagnostics.AddError(
			"login failed",
			"login failed: %s"+err.Error())
		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *freeipaProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewFreeipaHostResource,
	}
}

func (p *freeipaProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewFreeipaHostDataSource,
	}
}

func (p *freeipaProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &freeipaProvider{
			version: version,
		}
	}
}
