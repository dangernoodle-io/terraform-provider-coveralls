package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"terraform-provider-coveralls/internal/provider/client"
)

var (
	_ resource.Resource              = &RepositoryResource{}
	_ resource.ResourceWithConfigure = &RepositoryResource{}
)

func NewRepositoryResource() resource.Resource {
	return &RepositoryResource{}
}

type RepositoryResource struct {
	coveralls *Coveralls
}

func (r *RepositoryResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_repository"
}

func (r *RepositoryResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this resource to manage a Coveralls repository.",
		Attributes: map[string]schema.Attribute{
			"comment_on_pull_requests": schema.BoolAttribute{
				Description: "Whether comments should be posted on pull requests.",
				Required:    true,
			},
			"commit_status_fail_threshold": schema.Float64Attribute{
				Description: "Minimum coverage that must be present on a build for the build to pass.",
				Optional:    true,
			},
			"commit_status_fail_change_threshold": schema.Float64Attribute{
				Description: "Maximum allowed amount of decrease that will be allowed for the build to pass.",
				Optional:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "Date and time when the Coveralls repository was created.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"id": schema.StringAttribute{
				Description: "Unique identifier for the repository.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the repository in the form `<owner>/<name>`.",
				Required:            true,
			},
			"send_build_status": schema.BoolAttribute{
				Description: "Whether build status should be sent to the git provider.",
				Required:    true,
			},
			"service": schema.StringAttribute{
				MarkdownDescription: "Git provider, eg: `github`",
				Required:            true,
			},
			"token": schema.StringAttribute{
				Description: "Repository Token.",
				Computed:    true,
				Sensitive:   true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"updated_at": schema.StringAttribute{
				Description: "Date and time when the Coveralls repository was last updated.",
				Computed:    true,
			},
		},
	}
}

func (r *RepositoryResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	coveralls, ok := req.ProviderData.(*Coveralls)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *coveralls.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.coveralls = coveralls
}

func (r *RepositoryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan := &RepositoryState{}
	diags := req.Plan.Get(ctx, plan)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	repository := setRepositoryConfig(plan)
	// these are required on the struct during creation
	repository.Name = plan.Name.ValueString()
	repository.Service = plan.Service.ValueString()

	_, err := r.coveralls.client.Create(ctx, repository)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating repository",
			"Could not create repository, unexpected error: "+err.Error(),
		)
		return
	}

	// the 'token' isn't available on initial creation, so an additional read call is necessary
	repository, err = r.coveralls.client.Get(ctx, repository.Service, repository.Name)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading repository",
			"Could not read repository, unexpected error: "+err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, r.coveralls.converter(repository))
	resp.Diagnostics.Append(diags...)
}

func (r *RepositoryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	state := &RepositoryState{}
	diags := req.State.Get(ctx, state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := strings.Split(state.Id.ValueString(), ":")
	repository, err := r.coveralls.client.Get(ctx, id[0], id[1])

	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading repository",
			"Could not read repository, unexpected error: "+err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, r.coveralls.converter(repository))
	resp.Diagnostics.Append(diags...)
}

func (r *RepositoryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	plan := &RepositoryState{}
	diags := req.Plan.Get(ctx, plan)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := strings.Split(plan.Id.ValueString(), ":")
	repository := setRepositoryConfig(plan)

	_, err := r.coveralls.client.Update(ctx, id[0], id[1], repository)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating repository",
			"Could not update repository, unexpected error: "+err.Error(),
		)
		return
	}

	// the 'token' isn't available on initial creation, so an additional read call is necessary
	repository, err = r.coveralls.client.Get(ctx, id[0], id[1])

	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating repository",
			"Could not read repository, unexpected error: "+err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, r.coveralls.converter(repository))
	resp.Diagnostics.Append(diags...)
}

func (r *RepositoryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Warn(ctx, "Delete not supported by Coveralls API")
}

func (r *RepositoryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func setRepositoryConfig(state *RepositoryState) *client.Repository {
	return &client.Repository{
		CommentOnPullRequests: state.CommentOnPullRequests.ValueBool(),
		SendBuildStatus:       state.SendBuildStatus.ValueBool(),
		FailThreshold:         state.FailThreshold.ValueFloat64Pointer(),
		FailChangeThreshold:   state.FailChangeThreshold.ValueFloat64Pointer(),
	}
}
