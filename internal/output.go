package internal

type InvokeCloudBeesWorkflow struct {
	ComponentId      string            `json:"component_id,omitempty"`
	BranchName       string            `json:"branch_name,omitempty"`
	WorkflowFileName string            `json:"workflow_file_name,omitempty"`
	WorkflowInputs   map[string]string `json:"workflow_inputs,omitempty"`
}

type ProviderInfo struct {
	RunId      string `json:"run_id,omitempty"`
	RunAttempt string `json:"run_attempt,omitempty"`
	RunNumber  string `json:"run_number,omitempty"`
	JobName    string `json:"job_name,omitempty"`
	Provider   string `json:"provider,omitempty"`
}

type CloudEventData struct {
	ProviderInfo   ProviderInfo            `json:"provider_info,omitempty"`
	InvokeWorkflow InvokeCloudBeesWorkflow `json:"invoke_workflow,omitempty"`
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details []any  `json:"details"`
}

type SuccessResponse struct {
	Success      bool   `json:"success"`
	ErrorMessage string `json:"errorMessage"`
	EventOutput  struct {
		InvokeWorkflowOutput struct {
			RunUrl string `json:"runUrl"`
		} `json:"invokeWorkflowOutput"`
	} `json:"eventOutput"`
}
