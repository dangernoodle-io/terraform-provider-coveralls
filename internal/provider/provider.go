package provider

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-coveralls/internal/provider/client"
)

var (
	_ provider.Provider = &CoverallsProvider{}
)

type Coveralls struct {
	client    *client.Client
	converter RepositoryConverter
}

type CoverallsProvider struct {
	version string
}

type CoverallsProviderModel struct {
	Token types.String `tfsdk:"token"`
}

type RepositoryState struct {
	Id                    types.String  `tfsdk:"id"`
	Name                  types.String  `tfsdk:"name"`
	Service               types.String  `tfsdk:"service"`
	Token                 types.String  `tfsdk:"token"`
	CommentOnPullRequests types.Bool    `tfsdk:"comment_on_pull_requests"`
	SendBuildStatus       types.Bool    `tfsdk:"send_build_status"`
	FailThreshold         types.Float64 `tfsdk:"commit_status_fail_threshold"`
	FailChangeThreshold   types.Float64 `tfsdk:"commit_status_fail_change_threshold"`
	CreatedAt             types.String  `tfsdk:"created_at"`
	UpdatedAt             types.String  `tfsdk:"updated_at"`
}

type RepositoryConverter func(*client.Repository) *RepositoryState

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &CoverallsProvider{
			version: version,
		}
	}
}

func (p *CoverallsProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "coveralls"
	resp.Version = p.version
}

func (p *CoverallsProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"token": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

func (p *CoverallsProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config CoverallsProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Token.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Unknown token",
			"The provider cannot create the Client client")
	}

	if resp.Diagnostics.HasError() {
		return
	}

	token := os.Getenv("COVERALLS_API_TOKEN")

	if !config.Token.IsNull() {
		token = config.Token.ValueString()
	}

	if token == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("Token"),
			"Missing Client API token",
			"The provider cannot create the Client client as there is a missing API Token. "+
				"Set the token in the configuration or use the COVERALLS_API_TOKEN environment variable",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	c, err := client.NewCoveralls("https://coveralls.io", token)

	if err != nil {
		resp.Diagnostics.AddError("Error creating Client client", err.Error())
		return
	}

	coveralls := &Coveralls{
		client:    c,
		converter: repositoryConverter(),
	}

	resp.DataSourceData = coveralls
	resp.ResourceData = coveralls
}

func (p *CoverallsProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewRepositoryDataSource,
	}
}

func (p *CoverallsProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewRepositoryResource,
	}
}

func repositoryConverter() RepositoryConverter {
	return func(repository *client.Repository) *RepositoryState {
		return &RepositoryState{
			Id:                    types.StringValue(fmt.Sprintf("%s:%s", repository.Service, repository.Name)),
			Service:               types.StringValue(repository.Service),
			Name:                  types.StringValue(repository.Name),
			Token:                 types.StringValue(repository.Token),
			CommentOnPullRequests: types.BoolValue(repository.CommentOnPullRequests),
			SendBuildStatus:       types.BoolValue(repository.SendBuildStatus),
			FailThreshold:         types.Float64PointerValue(repository.FailThreshold),
			FailChangeThreshold:   types.Float64PointerValue(repository.FailChangeThreshold),
			CreatedAt:             types.StringValue(repository.CreatedAt),
			UpdatedAt:             types.StringValue(repository.UpdatedAt),
		}
	}
}
