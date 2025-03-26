package internal

import "context"

type Config struct {
	context.Context
	ComponentId       string            `json:"component-id,omitempty"`
	BranchName        string            `json:"branch-name,omitempty"`
	WorkflowFileName  string            `json:"workflow-file-name,omitempty"`
	WorkflowInputs    map[string]string `json:"workflow-inputs,omitempty"`
	GhaRunId          string            `json:"gha-run-id,omitempty"`
	GhaRunAttempt     string            `json:"gha-run-attempt,omitempty"`
	GhaRunNumber      string            `json:"gha-run-number,omitempty"`
	CloudBeesApiUrl   string            `json:"cloudbees-api-url,omitempty"`
	CloudBeesApiToken string            `json:"cloudbees-api-token,omitempty"`
	GhaRepository     string            `json:"gha-repository,omitempty"`
	GhaWorkflowRef    string            `json:"gha-workflow-ref,omitempty"`
	GhaServerUrl      string            `json:"gha-server-url,omitempty"`
	GhaJobName        string            `json:"gha-job-name,omitempty"`
}
