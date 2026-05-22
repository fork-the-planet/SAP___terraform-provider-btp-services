// btpservices/provider/cicd/jobs/resource_job_test.go

package cicdjobs_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/SAP/terraform-provider-sap-btp-services/btpservices/provider/cicd/utils"
)

func TestResourceCicdJob(t *testing.T) {
	t.Parallel()

	t.Run("error - missing name", func(t *testing.T) {
		t.Parallel()
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: utils.GetTestProviders(utils.Redacted, nil),
			Steps: []resource.TestStep{
				{
					Config: utils.HCLProviderBlock(utils.Redacted) + `
resource "btpservice_cicd_job" "test" {
  active               = true
  pipeline             = "sap-cloud-sdk"
  pipeline_version     = "3.0"
  pipeline_parameters  = "{}"
  build_retention_days = 28
  max_builds_to_keep   = 15
}
`,
					ExpectError: regexp.MustCompile(`The argument "name" is required`),
				},
			},
		})
	})

	t.Run("error - missing active", func(t *testing.T) {
		t.Parallel()
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: utils.GetTestProviders(utils.Redacted, nil),
			Steps: []resource.TestStep{
				{
					Config: utils.HCLProviderBlock(utils.Redacted) + `
resource "btpservice_cicd_job" "test" {
  name                 = "test-job"
  pipeline             = "sap-cloud-sdk"
  pipeline_version     = "3.0"
  pipeline_parameters  = "{}"
  build_retention_days = 28
  max_builds_to_keep   = 15
}
`,
					ExpectError: regexp.MustCompile(`The argument "active" is required`),
				},
			},
		})
	})

	t.Run("error - missing pipeline", func(t *testing.T) {
		t.Parallel()
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: utils.GetTestProviders(utils.Redacted, nil),
			Steps: []resource.TestStep{
				{
					Config: utils.HCLProviderBlock(utils.Redacted) + `
resource "btpservice_cicd_job" "test" {
  name                 = "test-job"
  active               = true
  pipeline_version     = "3.0"
  pipeline_parameters  = "{}"
  build_retention_days = 28
  max_builds_to_keep   = 15
}
`,
					ExpectError: regexp.MustCompile(`The argument "pipeline" is required`),
				},
			},
		})
	})

	t.Run("error - invalid pipeline value", func(t *testing.T) {
		t.Parallel()
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: utils.GetTestProviders(utils.Redacted, nil),
			Steps: []resource.TestStep{
				{
					Config: utils.HCLProviderBlock(utils.Redacted) + `
resource "btpservice_cicd_job" "test" {
  name                 = "test-job"
  active               = true
  pipeline             = "unknown-pipeline"
  pipeline_version     = "3.0"
  pipeline_parameters  = "{}"
  build_retention_days = 28
  max_builds_to_keep   = 15
}
`,
					ExpectError: regexp.MustCompile(`(?i)value must be one of`),
				},
			},
		})
	})

	t.Run("error - missing pipeline_version", func(t *testing.T) {
		t.Parallel()
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: utils.GetTestProviders(utils.Redacted, nil),
			Steps: []resource.TestStep{
				{
					Config: utils.HCLProviderBlock(utils.Redacted) + `
resource "btpservice_cicd_job" "test" {
  name                 = "test-job"
  active               = true
  pipeline             = "sap-cloud-sdk"
  pipeline_parameters  = "{}"
  build_retention_days = 28
  max_builds_to_keep   = 15
}
`,
					ExpectError: regexp.MustCompile(`The argument "pipeline_version" is required`),
				},
			},
		})
	})

	t.Run("error - missing pipeline_parameters", func(t *testing.T) {
		t.Parallel()
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: utils.GetTestProviders(utils.Redacted, nil),
			Steps: []resource.TestStep{
				{
					Config: utils.HCLProviderBlock(utils.Redacted) + `
resource "btpservice_cicd_job" "test" {
  name                 = "test-job"
  active               = true
  pipeline             = "sap-cloud-sdk"
  pipeline_version     = "3.0"
  build_retention_days = 28
  max_builds_to_keep   = 15
}
`,
					ExpectError: regexp.MustCompile(`The argument "pipeline_parameters" is required`),
				},
			},
		})
	})

	t.Run("error - invalid YAML in pipeline_parameters", func(t *testing.T) {
		t.Parallel()
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: utils.GetTestProviders(utils.Redacted, nil),
			Steps: []resource.TestStep{
				{
					Config: utils.HCLProviderBlock(utils.Redacted) + `
resource "btpservice_cicd_job" "test" {
  name                 = "test-job"
  active               = true
  pipeline             = "sap-cloud-sdk"
  pipeline_version     = "3.0"
  pipeline_parameters  = "key: [invalid yaml"
  build_retention_days = 28
  max_builds_to_keep   = 15
}
`,
					ExpectError: regexp.MustCompile(`(?i)invalid yaml`),
				},
			},
		})
	})

	t.Run("error - missing build_retention_days", func(t *testing.T) {
		t.Parallel()
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: utils.GetTestProviders(utils.Redacted, nil),
			Steps: []resource.TestStep{
				{
					Config: utils.HCLProviderBlock(utils.Redacted) + `
resource "btpservice_cicd_job" "test" {
  name                = "test-job"
  active              = true
  pipeline            = "sap-cloud-sdk"
  pipeline_version    = "3.0"
  pipeline_parameters = "{}"
  max_builds_to_keep  = 15
}
`,
					ExpectError: regexp.MustCompile(`The argument "build_retention_days" is required`),
				},
			},
		})
	})

	t.Run("error - missing max_builds_to_keep", func(t *testing.T) {
		t.Parallel()
		resource.Test(t, resource.TestCase{
			IsUnitTest:               true,
			ProtoV6ProviderFactories: utils.GetTestProviders(utils.Redacted, nil),
			Steps: []resource.TestStep{
				{
					Config: utils.HCLProviderBlock(utils.Redacted) + `
resource "btpservice_cicd_job" "test" {
  name                 = "test-job"
  active               = true
  pipeline             = "sap-cloud-sdk"
  pipeline_version     = "3.0"
  pipeline_parameters  = "{}"
  build_retention_days = 28
}
`,
					ExpectError: regexp.MustCompile(`The argument "max_builds_to_keep" is required`),
				},
			},
		})
	})
}
