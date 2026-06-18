// internal/cicd/client/facade_allowed_spaces.go

package cicdclient

import (
	"context"

	cicdmodels "github.com/SAP/terraform-provider-btp-services/internal/cicd/models"
)

type allowedSpacesFacade struct {
	hc *cicdHTTPClient
}

func newAllowedSpacesFacade(hc *cicdHTTPClient) allowedSpacesFacade {
	return allowedSpacesFacade{hc: hc}
}

// Get sends GET /v2/settings/allowedSpaces.
func (f *allowedSpacesFacade) Get(ctx context.Context) (*cicdmodels.AllowedSpacesResponse, error) {
	var result cicdmodels.AllowedSpacesResponse
	if err := f.hc.doGet(ctx, "/v2/settings/allowedSpaces", &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Set sends PUT /v2/settings/allowedSpaces, replacing the entire list.
// The API returns 204 with no body.
func (f *allowedSpacesFacade) Set(ctx context.Context, req cicdmodels.AllowedSpaceListDTO) error {
	return f.hc.doPut(ctx, "/v2/settings/allowedSpaces", req)
}
