---
page_title: "btpservice_cicd_job Resource - SAP BTP Services"
subcategory: ""
description: |-
  Manages a CI/CD job in the SAP BTP CI/CD service.
---

# btpservice_cicd_job (Resource)

Manages a CI/CD job in the SAP BTP CI/CD service.

## Example Usage

```terraform
# ---------------------------------------------------------------------------
# Prerequisite: a repository the jobs will use
# ---------------------------------------------------------------------------
resource "btpservice_cicd_repository" "app_repo" {
  name      = "my-app-repository"
  clone_url = "https://github.com/example/my-app"
}

# ---------------------------------------------------------------------------
# Credentials referenced by the pipeline parameters below
# ---------------------------------------------------------------------------
resource "btpservice_cicd_credential_basic_auth" "deploy_user" {
  name     = "cf-deploy-user"
  username = "deploy-user@example.com"
  password = var.cf_deploy_password
}

resource "btpservice_cicd_credential_kubernetes_config" "kubeconfig" {
  name        = "kube-config"
  kube_config = var.kube_config
}

resource "btpservice_cicd_credential_secret_text" "sonar_token" {
  name   = "sonar-token"
  secret = var.sonar_token
}

resource "btpservice_cicd_credential_container_registry" "registry" {
  name     = "container-registry"
  username = var.registry_username
  password = var.registry_password
  server   = "https://docker.io"
}

# =============================================================================
# Pipeline type: cf-env  (Cloud Foundry Environment)
# Full pipeline — all stages, runFirst/runLast, _additional credential and
# string variables, Cloud Transport Management, SonarCloud compliance.
# =============================================================================
resource "btpservice_cicd_job" "cf_env_full" {
  name                 = "cf-full-pipeline"
  description          = "CF environment pipeline with all stages"
  repository_id        = btpservice_cicd_repository.app_repo.id
  branch               = "main"
  pipeline             = "cf-env"
  pipeline_version     = "3.0"
  active               = true
  build_retention_days = 28
  max_builds_to_keep   = 10

  pipeline_parameters = <<-YAML
    configurationSource: job_parameter
    cfEnvConfiguration:
      stages:
        build:
          buildTool: mta
          buildToolVersion: MBTJ21N24
          buildDescriptor: java_app/mta.yaml
          npmLint:
            failOnError: true
          runFirst:
            command: echo "Starting build"
          runLast:
            command: echo "Build finished"
          _additional:
            stringVariables:
              - name: MAVEN_OPTS
                value: -Xmx2g
            credentialVariables:
              - name: NEXUS_PASSWORD
                valueSource: ${btpservice_cicd_credential_basic_auth.deploy_user.id}
        malwareScan:
          scan: true
        additionalTests:
          npmTests:
            npmScript: unit-test
            buildDescriptor: test/package.json
        acceptance:
          cfDeploy:
            strategy: default
            apiEndpoint: https://api.cf.us10.hana.ondemand.com
            org: my-org
            space: acceptance
            credential: ${btpservice_cicd_credential_basic_auth.deploy_user.id}
            mtaExtensionDescriptors:
              - mta-acc.mtaext
          webdriverIoTests:
            baseUrl: https://myapp-acc.cfapps.io
            npmScript: e2e
            credential: ${btpservice_cicd_credential_basic_auth.deploy_user.id}
            buildDescriptor: e2e/package.json
          runFirst:
            command: echo "Starting acceptance deploy"
          runLast:
            command: echo "Acceptance tests done"
          _additional:
            stringVariables:
              - name: TARGET_ENV
                value: acceptance
            credentialVariables:
              - name: ACC_API_KEY
                valueSource: ${btpservice_cicd_credential_basic_auth.deploy_user.id}
        compliance:
          sonarScan:
            mode: SonarCloud
            serverUrl: https://sonarcloud.io
            organization: my-org
            projectKey: my-org_my-project
            tokenCredential: ${btpservice_cicd_credential_secret_text.sonar_token.id}
        release:
          cfDeploy:
            strategy: blue-green
            apiEndpoint: https://api.cf.us10.hana.ondemand.com
            org: my-org
            space: production
            credential: ${btpservice_cicd_credential_basic_auth.deploy_user.id}
            mtaExtensionDescriptors:
              - mta-prod.mtaext
          cloudTransportManagement:
            nodeName: prod-ctm-node
            credential: ${btpservice_cicd_credential_basic_auth.deploy_user.id}
            nodeOperation: export
          runFirst:
            command: echo "Starting release"
          runLast:
            command: echo "Release complete"
          _additional:
            stringVariables:
              - name: TARGET_ENV
                value: production
            credentialVariables:
              - name: PROD_API_KEY
                valueSource: ${btpservice_cicd_credential_basic_auth.deploy_user.id}
      lifeCycle:
        afterAllStages:
          command: echo "Pipeline finished"
          credentialVariables: []
          stringVariables: []
  YAML
}

# =============================================================================
# Pipeline type: cf-env — pipeline_parameters loaded from a YAML file
# Use file() when the pipeline config is large, shared across jobs, or managed
# separately from the Terraform configuration.
#
# Create a file named cf_pipeline_params.yaml next to this .tf file with the
# full pipeline parameters YAML. Credential IDs must be hardcoded in the file
# because file() reads raw content with no variable substitution. Example:
#
#   configurationSource: job_parameter
#   cfEnvConfiguration:
#     stages:
#       build:
#         buildTool: mta
#         buildToolVersion: MBTJ21N24
#         runFirst:
#           command: echo "Starting build"
#         _additional:
#           stringVariables:
#             - name: MAVEN_OPTS
#               value: -Xmx2g
#           credentialVariables:
#             - name: NEXUS_PASSWORD
#               valueSource: <credential-id>
#       release:
#         cfDeploy:
#           strategy: blue-green
#           apiEndpoint: https://api.cf.us10.hana.ondemand.com
#           org: my-org
#           space: production
#           credential: <credential-id>
# =============================================================================
resource "btpservice_cicd_job" "cf_env_from_file" {
  name                 = "cf-from-file"
  description          = "Pipeline parameters loaded from an external YAML file"
  repository_id        = btpservice_cicd_repository.app_repo.id
  branch               = "main"
  pipeline             = "cf-env"
  pipeline_version     = "3.0"
  active               = true
  build_retention_days = 28
  max_builds_to_keep   = 10

  pipeline_parameters = file("${path.module}/cf_pipeline_params.yaml")
}

# =============================================================================
# Pipeline type: cf-env — pipeline_parameters rendered from a template file
# Use templatefile() to inject credential IDs (or any Terraform value) into
# the YAML at plan time, avoiding hardcoded IDs in the YAML file itself.
#
# Create a file named cf_pipeline_params.tftpl next to this .tf file using
# ${variable_name} placeholders for any value you want to inject. Example:
#
#   configurationSource: job_parameter
#   cfEnvConfiguration:
#     stages:
#       build:
#         buildTool: mta
#         buildToolVersion: MBTJ21N24
#         runFirst:
#           command: echo "Starting build"
#         _additional:
#           stringVariables:
#             - name: MAVEN_OPTS
#               value: -Xmx2g
#           credentialVariables:
#             - name: NEXUS_PASSWORD
#               valueSource: ${deploy_cred}
#       release:
#         cfDeploy:
#           strategy: blue-green
#           apiEndpoint: https://api.cf.us10.hana.ondemand.com
#           org: my-org
#           space: production
#           credential: ${deploy_cred}
#       compliance:
#         sonarScan:
#           mode: SonarCloud
#           serverUrl: https://sonarcloud.io
#           organization: my-org
#           projectKey: my-org_my-project
#           tokenCredential: ${sonar_cred}
# =============================================================================
resource "btpservice_cicd_job" "cf_env_from_template" {
  name                 = "cf-from-template"
  description          = "Pipeline parameters rendered from a template file via templatefile()"
  repository_id        = btpservice_cicd_repository.app_repo.id
  branch               = "main"
  pipeline             = "cf-env"
  pipeline_version     = "3.0"
  active               = true
  build_retention_days = 28
  max_builds_to_keep   = 10

  pipeline_parameters = templatefile("${path.module}/cf_pipeline_params.tftpl", {
    deploy_cred = btpservice_cicd_credential_basic_auth.deploy_user.id
    sonar_cred  = btpservice_cicd_credential_secret_text.sonar_token.id
  })
}

# =============================================================================
# Pipeline type: kyma-cnb  (Kyma Runtime — Cloud Native Buildpacks)
# Full pipeline — all stages, multiple images with exportHelmValues, helmValues
# overrides (literal + file source), runFirst/runLast, _additional credential
# and string variables across build, acceptance, compliance, and release.
# =============================================================================
resource "btpservice_cicd_job" "kyma_cnb_full" {
  name                 = "kyma-full-pipeline"
  description          = "Kyma CNB pipeline with all stages"
  repository_id        = btpservice_cicd_repository.app_repo.id
  branch               = "main"
  pipeline             = "kyma-cnb"
  pipeline_version     = "1.0"
  active               = true
  build_retention_days = 14
  max_builds_to_keep   = 20

  pipeline_parameters = <<-YAML
    configurationSource: job_parameter
    kymaCnbConfiguration:
      common:
        chartPath: helm/chart
        images:
          - name: backend-image
            path: backend
            tag: latest
            exportHelmValues:
              tagValuePath: backend.image.tag
              repositoryValuePath: backend.image.repository
          - name: frontend-image
            path: frontend
            tag: latest
            exportHelmValues:
              tagValuePath: frontend.image.tag
              repositoryValuePath: frontend.image.repository
      stages:
        build:
          npm:
            toolVersion: N24
            script: ci-build
          lint:
            active: true
            failOnError: true
          cnb:
            containerRegistry:
              url: https://docker.io/myuser
              credential: ${btpservice_cicd_credential_container_registry.registry.id}
          helmDependencyUpdate:
            active: true
            sourceRepositories:
              - name: my-helm-repo
                url: https://charts.example.com
          runFirst:
            command: npm ci --prefer-offline
          runLast:
            command: echo "Build stage finished"
          _additional:
            stringVariables:
              - name: NODE_ENV
                value: production
              - name: APP_VERSION
                value: "2.0.0"
            credentialVariables:
              - name: NPM_TOKEN
                valueSource: ${btpservice_cicd_credential_container_registry.registry.id}
        additionalTest:
          active: true
          npmExecuteScriptsRunScript: test
          _additional:
            credentialVariables: []
            stringVariables: []
        acceptance:
          active: true
          deploy:
            kubeConfigFileCredential: ${btpservice_cicd_credential_kubernetes_config.kubeconfig.id}
            namespace: staging
            helmReleaseName: my-app-staging
            helmValueFiles:
              - values-staging.yaml
            helmValues:
              - path: replicaCount
                value: "2"
                source: literal
              - path: resources.limits.memory
                value: 256Mi
                source: literal
              - path: config.secretFile
                value: secrets/staging.yaml
                source: file
          webdriverIoTest:
            active: true
            npmScript: wdi5
            baseUrl: https://my-app-staging.example.com
          runFirst:
            command: echo "Starting acceptance deploy"
          runLast:
            command: echo "Acceptance tests done"
          _additional:
            stringVariables:
              - name: TARGET_ENV
                value: staging
            credentialVariables:
              - name: STAGING_SECRET
                valueSource: ${btpservice_cicd_credential_kubernetes_config.kubeconfig.id}
        compliance:
          sonarExecuteScan:
            active: true
            mode: SonarCloud
            serverUrl: https://sonarcloud.io
            organization: my-org
            projectKey: my-org_my-project
            tokenCredential: ${btpservice_cicd_credential_secret_text.sonar_token.id}
        release:
          active: true
          deploy:
            kubeConfigFileCredential: ${btpservice_cicd_credential_kubernetes_config.kubeconfig.id}
            namespace: production
            helmReleaseName: my-app-prod
            helmValueFiles:
              - values-prod.yaml
            helmValues:
              - path: replicaCount
                value: "5"
                source: literal
              - path: resources.limits.cpu
                value: "1000m"
                source: literal
              - path: resources.limits.memory
                value: 512Mi
                source: literal
              - path: config.secretFile
                value: secrets/prod.yaml
                source: file
          runFirst:
            command: echo "Starting production release"
          runLast:
            command: echo "Release complete"
          _additional:
            stringVariables:
              - name: TARGET_ENV
                value: production
            credentialVariables:
              - name: PROD_SECRET
                valueSource: ${btpservice_cicd_credential_kubernetes_config.kubeconfig.id}
  YAML
}

# =============================================================================
# Pipeline type: cpi  (SAP Integration Suite Artifacts)
# Full pipeline — all stages active, all optional fields.
# =============================================================================
resource "btpservice_cicd_job" "cpi_full" {
  name                 = "cpi-full-pipeline"
  description          = "Integration Suite pipeline with all stages"
  repository_id        = btpservice_cicd_repository.app_repo.id
  branch               = "main"
  pipeline             = "cpi"
  pipeline_version     = "2.0"
  active               = true
  build_retention_days = 14
  max_builds_to_keep   = 20

  pipeline_parameters = <<-YAML
    configurationSource: job_parameter
    cpiConfiguration:
      common:
        integrationFlowId: MyIntegrationFlow
        apiPlanServiceKey: ${btpservice_cicd_credential_basic_auth.deploy_user.id}
      stages:
        upload:
          integrationFlowName: My Integration Flow Name
          packageId: com.example.mypackage
          active: true
        deploy:
          active: true
        integrationTests:
          active: true
          integrationFlowPlanServiceKey: ${btpservice_cicd_credential_basic_auth.deploy_user.id}
          contentType: application/json
          messageBodyPath: test/payload.json
  YAML
}

# =============================================================================
# notification_configuration — ANS (SAP Alert Notification Service)
# Attach an ANS credential to a job so build events are forwarded to the
# Alert Notification Service. The credential must already exist.
# =============================================================================
resource "btpservice_cicd_credential_basic_auth" "ans_credential" {
  name        = "ans-service-key"
  description = "ANS service key credential"
  username    = "ans-client-id"
  password    = var.ans_client_secret
}

resource "btpservice_cicd_job" "cf_with_ans" {
  name                 = "cf-pipeline-with-ans"
  description          = "CF pipeline that sends build notifications via ANS"
  repository_id        = btpservice_cicd_repository.app_repo.id
  branch               = "main"
  pipeline             = "cf-env"
  pipeline_version     = "3.0"
  active               = true
  build_retention_days = 28
  max_builds_to_keep   = 10

  pipeline_parameters = <<-YAML
    configurationSource: source_repository
  YAML

  notification_configuration = {
    ans = {
      active        = true
      credential_id = btpservice_cicd_credential_basic_auth.ans_credential.id
      custom_tag    = "my-team"
    }
  }
}

# Disable ANS notifications without removing the configuration block.
resource "btpservice_cicd_job" "cf_ans_disabled" {
  name                 = "cf-pipeline-ans-off"
  description          = "CF pipeline with ANS notifications disabled"
  repository_id        = btpservice_cicd_repository.app_repo.id
  branch               = "main"
  pipeline             = "cf-env"
  pipeline_version     = "3.0"
  active               = true
  build_retention_days = 28
  max_builds_to_keep   = 10

  pipeline_parameters = <<-YAML
    configurationSource: source_repository
  YAML

  notification_configuration = {
    ans = {
      active        = false
      credential_id = btpservice_cicd_credential_basic_auth.ans_credential.id
    }
  }
}

# =============================================================================
# Pipeline type: sap-ui5-abap-fes  (SAP UI5 for ABAP Platform / Fiori)
# Full pipeline — all stages, build tool version, lint, additional test,
# malware scan, SonarCloud compliance, release with transport request.
# =============================================================================
resource "btpservice_cicd_job" "ui5_full" {
  name                 = "ui5-full-pipeline"
  description          = "UI5 pipeline deploying to ABAP FES with all stages"
  repository_id        = btpservice_cicd_repository.app_repo.id
  branch               = "main"
  pipeline             = "sap-ui5-abap-fes"
  pipeline_version     = "1.0"
  active               = true
  build_retention_days = 14
  max_builds_to_keep   = 20

  pipeline_parameters = <<-YAML
    configurationSource: job_parameter
    sapUi5FesConfiguration:
      build:
        buildToolVersion: N24
        npmExecuteLint:
          active: true
          failOnError: true
        runFirst:
          command: npm install
        runLast:
          command: npm run build
        _additional:
          stringVariables:
            - name: NODE_ENV
              value: production
          credentialVariables:
            - name: NPM_TOKEN
              valueSource: ${btpservice_cicd_credential_container_registry.registry.id}
      additionalTest:
        active: true
        npmExecuteScriptsRunScript: test
        _additional:
          credentialVariables: []
          stringVariables: []
      runMalwareScan: true
      compliance:
        sonarExecuteScan:
          active: true
          mode: SonarCloud
          serverUrl: https://sonarcloud.io
          organization: my-org
          projectKey: my-ui5-project-key
          tokenCredential: ${btpservice_cicd_credential_secret_text.sonar_token.id}
      release:
        active: true
        abapEndpoint: https://your-abap-system.com
        abapPackage: Z_UI5_RESOURCES
        applicationName: Z_CUSTOM_APP
        applicationDescription: Deployment via CI/CD Service
        uploadCredential: ${btpservice_cicd_credential_basic_auth.deploy_user.id}
        transportRequestIdSource: parameter
        transportRequestId: S4HK900001
        _additional:
          credentialVariables: []
          stringVariables: []
  YAML
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `active` (Boolean) Whether the job is active. Inactive jobs cannot be executed.
- `branch` (String) Branch pattern for the job. Required when `repository_id` is set.
- `build_retention_days` (Number) Number of days build artifacts are retained. Must be between 1 and 28 (inclusive).
- `max_builds_to_keep` (Number) Maximum number of builds retained for this job.
- `name` (String) Name of the job. Must match `[a-zA-Z0-9_-]{1,64}`.
- `pipeline` (String) Pipeline type. One of: `cpi`, `cf-env`, `kyma-cnb`, `sap-ui5-abap-fes`.
- `pipeline_parameters` (String) Pipeline parameters as a YAML string. Use `file()` or `templatefile()` to load from a file. When `configurationSource` is `source_repository`, the pipeline reads its config from the repo — set this to `configurationSource: source_repository`. When `configurationSource` is `job_parameter`, provide the full pipeline configuration here. The value is stored as-is in state so formatting is preserved across plans.
- `pipeline_version` (String) Version of the pipeline type (e.g. `3.0`, `1.0`).
- `repository_id` (String) ID of the source repository used by this job.

### Optional

- `description` (String) Optional human-readable description of the job.
- `notification_configuration` (Attributes) Optional notification settings for the job. (see [below for nested schema](#nestedatt--notification_configuration))

### Read-Only

- `id` (String) Unique identifier of the job (assigned by the API).

<a id="nestedatt--notification_configuration"></a>
### Nested Schema for `notification_configuration`

Optional:

- `ans` (Attributes) SAP Alert Notification Service (ANS) settings. (see [below for nested schema](#nestedatt--notification_configuration--ans))

<a id="nestedatt--notification_configuration--ans"></a>
### Nested Schema for `notification_configuration.ans`

Required:

- `active` (Boolean) Whether ANS notifications are active for this job.
- `credential_id` (String) ID of the ANS credential to use.

Optional:

- `custom_tag` (String) Optional custom tag added to ANS notifications.

## Import

Import is supported using the following syntax:

```terraform
# terraform import btpservice_cicd_job.<resource_name> <id>

terraform import btpservice_cicd_job.example pb091fd5-845b-4146-9bfc-d8cb74be04f8

# terraform import using id attribute in import block

import {
  to = btpservice_cicd_job.<resource_name>
  id = "<id>"
}

import {
  to =  btpservice_cicd_job.<resource_name>
  identity = {
   id = "<id>"
  }
}
```
