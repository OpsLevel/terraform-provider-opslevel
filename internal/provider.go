package provider

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/opslevel/opslevel-go/v2024"
)

// Ensure the implementation satisfies the provider.Provider interface.
var _ provider.Provider = &OpslevelProvider{}

type OpslevelProvider struct {
	version string
}

type OpslevelProviderModel struct {
	ApiToken   types.String `tfsdk:"api_token"`
	ApiUrl     types.String `tfsdk:"api_url"`
	ApiTimeout types.Int64  `tfsdk:"api_timeout"`
}

func (p *OpslevelProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "opslevel"
	resp.Version = p.version
}

func (p *OpslevelProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_token": schema.StringAttribute{
				Required:    true,
				Description: "The API authorization token. It can also be sourced from the OPSLEVEL_API_TOKEN environment variable.",
				Sensitive:   true,
			},
			"api_url": schema.StringAttribute{
				Optional:    true,
				Description: "The url of the OpsLevel API to. It can also be sourced from the OPSLEVEL_API_URL environment variable.",
				Sensitive:   false,
			},
			"api_timeout": schema.Int64Attribute{
				Optional:    true,
				Description: "Value (in seconds) to use for the timeout of API calls made",
				Sensitive:   false,
			},
		},
	}
}

func (p *OpslevelProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data OpslevelProviderModel
	defaultApiUrl := "https://api.opslevel.com/"
	defaultApiTimeout := 10
	tflog.Info(ctx, "Initializing opslevel client")

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values are now available.
	tflog.Debug(ctx, "Setting opslevel client API token...")
	if data.ApiToken.ValueString() == "" {
		if apiToken, ok := os.LookupEnv("OPSLEVEL_API_TOKEN"); ok {
			data.ApiToken = types.StringValue(apiToken)
		} else {
			resp.Diagnostics.AddError(
				"Missing OPSLEVEL_API_TOKEN",
				"An OPSLEVEL_API_TOKEN is needed to authenticate with the opslevel client. "+
					"This can be set as an environment variable or in the provider configuration block as 'api_token'.",
			)
		}
	}
	tflog.Debug(ctx, "opslevel client API token is set")

	tflog.Debug(ctx, "Setting opslevel client API endpoint URL...")
	if data.ApiUrl.IsNull() {
		if apiUrl, ok := os.LookupEnv("OPSLEVEL_API_URL"); ok {
			data.ApiUrl = types.StringValue(apiUrl)
		} else {
			data.ApiUrl = types.StringValue(defaultApiUrl)
		}
	}
	tflog.Debug(ctx, "opslevel client API endpoint URL is set")

	tflog.Debug(ctx, "Setting opslevel client API timeout...")
	if data.ApiTimeout.IsNull() {
		data.ApiTimeout = types.Int64Value(int64(defaultApiTimeout))

		if apiTimeout, ok := os.LookupEnv("OPSLEVEL_API_TIMEOUT"); ok {
			if timeout, err := strconv.Atoi(apiTimeout); err == nil {
				data.ApiTimeout = types.Int64Value(int64(timeout))
			} else {
				resp.Diagnostics.AddWarning(
					"Expected OPSLEVEL_API_TIMEOUT to be an int",
					fmt.Sprintf("OPSLEVEL_API_TIMEOUT was set but not as an int. The default timeout value of %d seconds will be used.", defaultApiTimeout),
				)
			}
		}
	}
	tflog.Debug(ctx, "opslevel client API timeout is set")

	opts := []opslevel.Option{
		opslevel.SetAPIToken(data.ApiToken.ValueString()),
		opslevel.SetURL(data.ApiUrl.ValueString()),
		opslevel.SetTimeout(time.Second * time.Duration(data.ApiTimeout.ValueInt64())),
		opslevel.SetUserAgentExtra(fmt.Sprintf("terraform-provider-%s", p.version)),
	}
	client := opslevel.NewGQLClient(opts...)

	tflog.Debug(ctx, "Validating OpsLevel client...")
	if err := client.Validate(); err != nil {
		tflog.Error(ctx, fmt.Sprintf("OpsLevel client validation error: %s", err))
	}
	tflog.Debug(ctx, "OpsLevel client is valid")
	tflog.Info(ctx, "OpsLevel client is initialized")

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *OpslevelProvider) Resources(context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewDomainResource,
	}
}

func (p *OpslevelProvider) DataSources(context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewDomainDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &OpslevelProvider{version: version}
	}
}
