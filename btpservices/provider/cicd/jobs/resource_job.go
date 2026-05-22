// btpservices/provider/cicd/jobs/resource_job.go

package cicdjobs

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"gopkg.in/yaml.v3"

	cicdclient "github.com/SAP/terraform-provider-sap-btp-services/internal/cicd/client"
	cicdmodels "github.com/SAP/terraform-provider-sap-btp-services/internal/cicd/models"
	"github.com/SAP/terraform-provider-sap-btp-services/internal/shared"
)

var _ resource.Resource = &jobResource{}
var _ resource.ResourceWithConfigure = &jobResource{}

func NewJobResource() resource.Resource {
	return &jobResource{}
}

type jobResource struct {
	cli *cicdclient.CicdClientFacade
}

func (r *jobResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_cicd_job", req.ProviderTypeName)
}

func (r *jobResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a CI/CD job in the SAP BTP CI/CD service.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier of the job (assigned by the API).",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the job. Must match `[a-zA-Z0-9_-]{1,64}`.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Optional human-readable description of the job.",
				Optional:            true,
				Computed:            true,
			},
			"active": schema.BoolAttribute{
				MarkdownDescription: "Whether the job is active. Inactive jobs cannot be executed.",
				Required:            true,
			},
			"pipeline": schema.StringAttribute{
				MarkdownDescription: "Pipeline type. Known values: `sap-cloud-sdk`, `cpi`, `cf-env`, `kyma-cnb`, `sap-ui5-abap-fes`.",
				Required:            true,
			},
			"pipeline_version": schema.StringAttribute{
				MarkdownDescription: "Version of the pipeline type (e.g. `3.0`, `1.0`).",
				Required:            true,
			},
			"pipeline_parameters": schema.StringAttribute{
				MarkdownDescription: "Pipeline parameters as a YAML string. Use `file()` to load from a file. " +
					"The content is converted to JSON before being sent to the API.",
				Required: true,
				Validators: []validator.String{
					validYAMLValidator{},
				},
			},
			"build_retention_days": schema.Int64Attribute{
				MarkdownDescription: "Number of days build artifacts are retained.",
				Required:            true,
			},
			"max_builds_to_keep": schema.Int64Attribute{
				MarkdownDescription: "Maximum number of builds retained for this job.",
				Required:            true,
			},
			"branch": schema.StringAttribute{
				MarkdownDescription: "Branch pattern for the job. Required when `repository_id` is set.",
				Optional:            true,
				Computed:            true,
			},
			"repository_id": schema.StringAttribute{
				MarkdownDescription: "ID of the source repository used by this job.",
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

func (r *jobResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	clients, ok := req.ProviderData.(*shared.ProviderClients)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Provider Data Type",
			fmt.Sprintf("Expected *shared.ProviderClients, got: %T", req.ProviderData),
		)
		return
	}
	if clients.Cicd == nil {
		resp.Diagnostics.AddError(
			"Missing CI/CD Configuration",
			"A cicd{} block must be configured in the provider to use CI/CD resources.",
		)
		return
	}
	r.cli = clients.Cicd
}

func (r *jobResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan jobResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq, err := plan.toCreateRequest()
	if err != nil {
		resp.Diagnostics.AddError("Invalid pipeline_parameters", err.Error())
		return
	}

	if err := r.cli.Jobs.Create(ctx, createReq); err != nil {
		resp.Diagnostics.AddError("Error Creating Job", err.Error())
		return
	}

	result, err := r.cli.Jobs.Get(ctx, plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error Reading Job After Create", err.Error())
		return
	}

	state, err := jobResourceValueFrom(*result, plan.PipelineParameters.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error Mapping Job Response", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *jobResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state jobResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.cli.Jobs.Get(ctx, state.ID.ValueString())
	if err != nil {
		if cicdmodels.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error Reading Job", err.Error())
		return
	}

	updated, err := jobResourceValueFrom(*result, state.PipelineParameters.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error Mapping Job Response", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, updated)...)
}

func (r *jobResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan jobResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state jobResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateReq, err := plan.toUpdateRequest()
	if err != nil {
		resp.Diagnostics.AddError("Invalid pipeline_parameters", err.Error())
		return
	}

	if err := r.cli.Jobs.Update(ctx, updateReq); err != nil {
		resp.Diagnostics.AddError("Error Updating Job", err.Error())
		return
	}

	result, err := r.cli.Jobs.Get(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error Reading Job After Update", err.Error())
		return
	}

	updated, err := jobResourceValueFrom(*result, plan.PipelineParameters.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error Mapping Job Response", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, updated)...)
}

func (r *jobResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state jobResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.cli.Jobs.Delete(ctx, state.ID.ValueString()); err != nil {
		if !cicdmodels.IsNotFound(err) {
			resp.Diagnostics.AddError("Error Deleting Job", err.Error())
		}
	}
}

// validYAMLValidator is a schema validator that ensures a string is valid YAML.
type validYAMLValidator struct{}

func (v validYAMLValidator) Description(_ context.Context) string {
	return "value must be valid YAML"
}

func (v validYAMLValidator) MarkdownDescription(_ context.Context) string {
	return "value must be valid YAML"
}

func (v validYAMLValidator) ValidateString(_ context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	var out any
	if err := yaml.Unmarshal([]byte(req.ConfigValue.ValueString()), &out); err != nil {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid YAML",
			fmt.Sprintf("pipeline_parameters must be valid YAML: %s", err.Error()),
		)
	}
	_ = types.StringValue // keep import used
}
