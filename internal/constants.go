package internal

const (
	ComponentId      = "COMPONENT_ID"
	BranchName       = "BRANCH_NAME"
	WorkflowFileName = "WORKFLOW_FILE_NAME"
	WorkflowInputs   = "WORKFLOW_INPUTS"

	GithubRunId      = "GITHUB_RUN_ID"
	GithubRunAttempt = "GITHUB_RUN_ATTEMPT"
	GithubRunNumber  = "GITHUB_RUN_NUMBER"

	CloudbeesApiUrl              = "CLOUDBEES_API_URL"
	CloudbeesApiToken            = "CLOUDBEES_API_TOKEN"
	GithubRepository             = "GITHUB_REPOSITORY"
	GithubWorkflowRef            = "GITHUB_WORKFLOW_REF"
	GithubServerUrl              = "GITHUB_SERVER_URL"
	InvokeCloudBeesWorkflowEvent = "cloudbees.platform.invoke.workflow"
	SpecVersion                  = "1.0"
	ContentTypeJson              = "application/json"
	ContentTypeHeaderKey         = "Content-Type"
	ContentTypeCloudEventsJson   = "application/cloudevents+json"
	AuthorizationHeaderKey       = "Authorization"
	Bearer                       = "Bearer "
	PostMethod                   = "POST"
	GithubProvider               = "GITHUB"
	GithubJobName                = "GITHUB_JOB"
	GitHubRef                    = "GITHUB_REF"
)
