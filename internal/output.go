package internal

type InvokeCloudBeesWorkflow struct {
	ComponentId      string            `json:"component-id,omitempty"`
	BranchName       string            `json:"branch-name,omitempty"`
	WorkflowFileName string            `json:"workflow-file-name,omitempty"`
	WorkflowInputs   map[string]string `json:"workflow-inputs,omitempty"`
}

type ProviderInfo struct {
	RunId      string `json:"run_id,omitempty"`
	RunAttempt string `json:"run_attempt,omitempty"`
	RunNumber  string `json:"run_number,omitempty"`
	JobName    string `json:"job_name,omitempty"`
	Provider   string `json:"provider,omitempty"`
}

type Output struct {
	ProviderInfo   ProviderInfo            `json:"provider_info,omitempty"`
	InvokeWorkflow InvokeCloudBeesWorkflow `json:"invoke_workflow,omitempty"`
}
