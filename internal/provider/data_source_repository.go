package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ datasource.DataSource = &RepositoryDataSource{}

type RepositoryDataSource struct {
	coveralls *Coveralls
}

func NewRepositoryDataSource() datasource.DataSource {
	return &RepositoryDataSource{}
}

func (d *RepositoryDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_repository"
}

func (d *RepositoryDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to retrieve information about a Coveralls repository.",
		Attributes: map[string]schema.Attribute{
			"comment_on_pull_requests": schema.BoolAttribute{
				Description: "Whether comments should be posted on pull requests.",
				Computed:    true,
			},
			"commit_status_fail_threshold": schema.Float64Attribute{
				Description: "Minimum coverage that must be present on a build for the build to pass.",
				Computed:    true,
			},
			"commit_status_fail_change_threshold": schema.Float64Attribute{
				Description: "Maximum allowed amount of decrease that will be allowed for the build to pass.",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "Date and time when the Coveralls repository was created.",
				Computed:    true,
			},
			"id": schema.StringAttribute{
				Description: "Unique identifier for the repository.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the repository in the form `<owner>/<name>`.",
				Required:            true,
			},
			"send_build_status": schema.BoolAttribute{
				Description: "Whether build status should be sent to the git provider.",
				Computed:    true,
			},
			"service": schema.StringAttribute{
				MarkdownDescription: "Git provider, eg: `github`",
				Required:            true,
			},
			"token": schema.StringAttribute{
				Description: "Repository Token.",
				Computed:    true,
				Sensitive:   true,
			},
			"updated_at": schema.StringAttribute{
				Description: "Date and time when the Coveralls repository was last updated.",
				Computed:    true,
			},
		},
	}
}

func (d *RepositoryDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	coveralls, ok := req.ProviderData.(*Coveralls)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *coveralls.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.coveralls = coveralls
}

func (d *RepositoryDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	state := &RepositoryState{}
	resp.Diagnostics.Append(req.Config.Get(ctx, state)...)

	repository, err := d.coveralls.client.Get(ctx, state.Service.ValueString(), state.Name.ValueString())

	if err != nil {
		ctx = tflog.SetField(ctx, "error", err.Error())
		tflog.Error(ctx, "failed")

		resp.Diagnostics.AddError(
			"Unable to read repository data",
			"Could not read repository, unexpected error: "+err.Error(),
		)
		return
	}

	diags := resp.State.Set(ctx, d.coveralls.converter(repository))
	resp.Diagnostics.Append(diags...)
}
