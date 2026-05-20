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
