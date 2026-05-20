// internal/cicd/models/job.go

package cicdmodels

// Job is the API response model for a single CI/CD job.
type Job struct {
	ID                        string                     `json:"id"`
	APIVersion                string                     `json:"apiVersion"`
	Name                      string                     `json:"name"`
	Active                    bool                       `json:"active"`
	Description               string                     `json:"description,omitempty"`
	Pipeline                  string                     `json:"pipeline"`
	PipelineVersion           string                     `json:"pipelineVersion"`
	PipelineParameters        map[string]any             `json:"pipelineParameters"`
	BuildRetentionDays        int                        `json:"buildRetentionDays"`
	MaxBuildsToKeep           int                        `json:"maxBuildsToKeep"`
	Branch                    string                     `json:"branch,omitempty"`
	RepositoryID              string                     `json:"repositoryId,omitempty"`
	NotificationConfiguration *NotificationConfiguration `json:"notificationConfiguration,omitempty"`
}

// NotificationConfiguration holds the notification settings for a job.
type NotificationConfiguration struct {
	ANS *AnsConfiguration `json:"ans,omitempty"`
}

// AnsConfiguration holds the SAP Alert Notification Service settings for a job.
type AnsConfiguration struct {
	Active      bool   `json:"active"`
	CredentialID string `json:"credentialId"`
	CustomTag   string `json:"customTag,omitempty"`
}

// JobListResponse is the envelope returned by GET /v2/repositories/{reference}/jobs.
type JobListResponse struct {
	Embedded *JobListEmbedded `json:"_embedded,omitempty"`
}

// JobListEmbedded holds the nested jobs array.
type JobListEmbedded struct {
	Jobs []Job `json:"jobs"`
}
