// internal/cicd/models/allowed_spaces.go

package cicdmodels

// AllowedSpace is one entry in the allowed-spaces list.
// Both fields are required by the API (POST body and GET response).
type AllowedSpace struct {
	SpaceGUID string `json:"spaceGuid"`
	Comment   string `json:"comment"`
}

// AllowedSpaceListDTO is the PUT request body for /v2/settings/allowedSpaces.
// It replaces the entire list atomically.
type AllowedSpaceListDTO struct {
	AllowedSpaces []AllowedSpace `json:"allowedSpaces"`
}

// AllowedSpacesResponse is the GET response body from /v2/settings/allowedSpaces.
type AllowedSpacesResponse struct {
	AllowedSpaces []AllowedSpace `json:"allowedSpaces"`
}
