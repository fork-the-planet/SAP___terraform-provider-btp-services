// internal/cicd/client/facade_jobs.go

package cicdclient

import (
	"context"
	"fmt"
	"net/url"

	cicdmodels "github.com/SAP/terraform-provider-sap-btp-services/internal/cicd/models"
)

type jobsFacade struct {
	hc *cicdHTTPClient
}

func newJobsFacade(hc *cicdHTTPClient) jobsFacade {
	return jobsFacade{hc: hc}
}

// Create sends POST /v2/jobs.
// The API returns 201 with no body; a subsequent Get is required to read the assigned ID.
func (f *jobsFacade) Create(ctx context.Context, req cicdmodels.CreateJobRequest) error {
	return f.hc.doPost(ctx, "/v2/jobs", req)
}

// Get sends GET /v2/jobs/{reference} where reference is the job name or ID.
func (f *jobsFacade) Get(ctx context.Context, reference string) (*cicdmodels.Job, error) {
	var result cicdmodels.Job
	err := f.hc.doGet(ctx, fmt.Sprintf("/v2/jobs/%s", url.PathEscape(reference)), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Update sends PUT /v2/jobs (full replace, job identified by name in request body).
// The API returns 204 with no body.
func (f *jobsFacade) Update(ctx context.Context, req cicdmodels.UpdateJobRequest) error {
	return f.hc.doPut(ctx, "/v2/jobs", req)
}

// Delete sends DELETE /v2/jobs/{reference}.
// The API returns 204 with no body.
func (f *jobsFacade) Delete(ctx context.Context, reference string) error {
	return f.hc.doDelete(ctx, fmt.Sprintf("/v2/jobs/%s", url.PathEscape(reference)))
}

// ListByRepository sends GET /v2/repositories/{reference}/jobs and returns all jobs for the repository.
// reference is the repository name or ID.
func (f *jobsFacade) ListByRepository(ctx context.Context, repositoryReference string) ([]cicdmodels.Job, error) {
	var result cicdmodels.JobListResponse
	err := f.hc.doGet(ctx, fmt.Sprintf("/v2/repositories/%s/jobs", url.PathEscape(repositoryReference)), &result)
	if err != nil {
		return nil, err
	}
	if result.Embedded == nil {
		return []cicdmodels.Job{}, nil
	}
	return result.Embedded.Jobs, nil
}
