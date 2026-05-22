// btpservices/provider/cicd/jobs/types.go

package cicdjobs

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"gopkg.in/yaml.v3"

	cicdmodels "github.com/SAP/terraform-provider-sap-btp-services/internal/cicd/models"
)

// jobResourceModel is the Terraform state model for the job resource.
type jobResourceModel struct {
	ID                 types.String `tfsdk:"id"`
	Name               types.String `tfsdk:"name"`
	Description        types.String `tfsdk:"description"`
	Active             types.Bool   `tfsdk:"active"`
	Pipeline           types.String `tfsdk:"pipeline"`
	PipelineVersion    types.String `tfsdk:"pipeline_version"`
	PipelineParameters types.String `tfsdk:"pipeline_parameters"`
	BuildRetentionDays types.Int64  `tfsdk:"build_retention_days"`
	MaxBuildsToKeep    types.Int64  `tfsdk:"max_builds_to_keep"`
	Branch             types.String `tfsdk:"branch"`
	RepositoryID       types.String `tfsdk:"repository_id"`
}

// yamlToMap parses a YAML string into map[string]any for the API request.
func yamlToMap(yamlStr string) (map[string]any, error) {
	var m map[string]any
	if err := yaml.Unmarshal([]byte(yamlStr), &m); err != nil {
		return nil, fmt.Errorf("invalid YAML in pipeline_parameters: %w", err)
	}
	if m == nil {
		m = map[string]any{}
	}
	return m, nil
}

func jobResourceValueFrom(v cicdmodels.Job, priorYAML string) (jobResourceModel, error) {
	m := jobResourceModel{
		ID:                 types.StringValue(v.ID),
		Name:               types.StringValue(v.Name),
		Description:        types.StringValue(v.Description),
		Active:             types.BoolValue(v.Active),
		Pipeline:           types.StringValue(v.Pipeline),
		PipelineVersion:    types.StringValue(v.PipelineVersion),
		BuildRetentionDays: types.Int64Value(int64(v.BuildRetentionDays)),
		MaxBuildsToKeep:    types.Int64Value(int64(v.MaxBuildsToKeep)),
		Branch:             types.StringValue(v.Branch),
		RepositoryID:       types.StringValue(v.RepositoryID),
	}

	// The API re-serializes pipelineParameters with alphabetically sorted keys and
	// different indentation — comparing it byte-for-byte with the user's YAML would
	// always produce a diff. Always preserve the prior YAML string so Terraform sees
	// no change unless the user actually edits the attribute.
	m.PipelineParameters = types.StringValue(priorYAML)

	return m, nil
}

func (m jobResourceModel) toCreateRequest() (cicdmodels.CreateJobRequest, error) {
	params, err := yamlToMap(m.PipelineParameters.ValueString())
	if err != nil {
		return cicdmodels.CreateJobRequest{}, err
	}
	return cicdmodels.CreateJobRequest{
		Name:               m.Name.ValueString(),
		Description:        m.Description.ValueString(),
		Active:             m.Active.ValueBool(),
		Pipeline:           m.Pipeline.ValueString(),
		PipelineVersion:    m.PipelineVersion.ValueString(),
		PipelineParameters: params,
		BuildRetentionDays: m.BuildRetentionDays.ValueInt64(),
		MaxBuildsToKeep:    m.MaxBuildsToKeep.ValueInt64(),
		Branch:             m.Branch.ValueString(),
		RepositoryID:       m.RepositoryID.ValueString(),
	}, nil
}

func (m jobResourceModel) toUpdateRequest() (cicdmodels.UpdateJobRequest, error) {
	params, err := yamlToMap(m.PipelineParameters.ValueString())
	if err != nil {
		return cicdmodels.UpdateJobRequest{}, err
	}
	return cicdmodels.UpdateJobRequest{
		ID:                 m.ID.ValueString(),
		Name:               m.Name.ValueString(),
		Description:        m.Description.ValueString(),
		Active:             m.Active.ValueBool(),
		Pipeline:           m.Pipeline.ValueString(),
		PipelineVersion:    m.PipelineVersion.ValueString(),
		PipelineParameters: params,
		BuildRetentionDays: m.BuildRetentionDays.ValueInt64(),
		MaxBuildsToKeep:    m.MaxBuildsToKeep.ValueInt64(),
		Branch:             m.Branch.ValueString(),
		RepositoryID:       m.RepositoryID.ValueString(),
	}, nil
}
