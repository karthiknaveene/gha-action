package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func (config *Config) Run(_ context.Context) (err error) {

	validationError := setEnvVars(config)
	if validationError != nil {
		return validationError
	}

	cloudEventData := prepareCloudEventData(config)

	cloudEvent, err := prepareCloudEvent(config, cloudEventData)
	if err != nil {
		fmt.Println("error preparing CloudEvent %s", err)
		return nil
	}
	err = sendCloudEvent(cloudEvent, config)
	if err != nil {
		fmt.Println("error sending CloudEvent %s", err)
		return nil
	}
	return nil
}

func setEnvVars(cfg *Config) error {
	ghaRunId := os.Getenv(GithubRunId)
	if ghaRunId == "" {
		return fmt.Errorf(GithubRunId + " is not set in the environment")
	}
	cfg.GhaRunId = ghaRunId

	ghaRunAttempt := os.Getenv(GithubRunAttempt)
	if ghaRunAttempt == "" {
		return fmt.Errorf(GithubRunAttempt + " is not set in the environment")
	}
	cfg.GhaRunAttempt = ghaRunAttempt

	cloudBeesApiUrl := os.Getenv(CloudbeesApiUrl)
	if cloudBeesApiUrl == "" {
		return fmt.Errorf(CloudbeesApiUrl + " is not set in the environment")
	}
	cfg.CloudBeesApiUrl = cloudBeesApiUrl

	cloudBeesApiToken := os.Getenv(CloudbeesApiToken)
	if cloudBeesApiToken == "" {
		return fmt.Errorf(CloudbeesApiToken + " is not set in the environment")
	}
	cfg.CloudBeesApiToken = cloudBeesApiToken

	ghaRunNumber := os.Getenv(GithubRunNumber)
	if ghaRunNumber == "" {
		return fmt.Errorf(GithubRunNumber + " is not set in the environment")
	}

	cfg.GhaRunNumber = ghaRunNumber

	ghaRepository := os.Getenv(GithubRepository)
	if ghaRepository == "" {
		return fmt.Errorf(GithubRepository + " is not set in the environment")
	}

	cfg.GhaRepository = ghaRepository

	ghaWorkflowRef := os.Getenv(GithubWorkflowRef)
	if ghaWorkflowRef == "" {
		return fmt.Errorf(GithubWorkflowRef + " is not set in the environment")
	}

	cfg.GhaWorkflowRef = ghaWorkflowRef

	ghaJobName := os.Getenv(GithubJobName)
	if ghaJobName == "" {
		return fmt.Errorf(GithubJobName + " is not set in the environment")
	}

	cfg.GhaJobName = ghaJobName

	workflowName := os.Getenv(WorkflowFileName)
	if workflowName == "" {
		return fmt.Errorf(WorkflowFileName + " is not set in the environment")
	}
	cfg.WorkflowFileName = workflowName

	branchName := os.Getenv(BranchName)
	if branchName == "" {
		var err error
		branchName, err = getCurrentBranchFromRef()
		if err != nil {
			return fmt.Errorf(BranchName + " is not set in the environment")
		}
	}
	cfg.BranchName = branchName

	cfg.GhaServerUrl = os.Getenv(GithubServerUrl)

	return nil
}

func getCloudbeesFullUrl(config *Config) string {
	if !strings.HasSuffix(config.CloudBeesApiUrl, "/") {
		config.CloudBeesApiUrl += "/"
	}
	return config.CloudBeesApiUrl + "v3/external-events"
}

func getSubject(config *Config) string {
	return config.GhaWorkflowRef + "|" + config.GhaRunId + "|" + config.GhaRunAttempt + "|" + config.GhaRunNumber
}

func getSource(config *Config) string {
	sourcePrefix := GithubProvider
	if config.GhaServerUrl != "" {
		sourcePrefix = config.GhaServerUrl + "/"
	}
	return sourcePrefix + config.GhaRepository
}
func prepareCloudEvent(config *Config, output Output) (cloudevents.Event, error) {
	cloudEvent := cloudevents.NewEvent()
	cloudEvent.SetID(uuid.NewString())
	cloudEvent.SetSubject(getSubject(config))
	cloudEvent.SetType(InvokeCloudBeesWorkflowEvent)
	cloudEvent.SetSource(getSource(config))
	cloudEvent.SetSpecVersion(SpecVersion)
	cloudEvent.SetTime(time.Now())
	err := cloudEvent.SetData(ContentTypeJson, output)
	fmt.Println("CloudEvent set data")
	fmt.Println(PrettyPrint(cloudEvent))
	if err != nil {
		return cloudevents.Event{}, fmt.Errorf("failed to set data: %v", err)
	}
	return cloudEvent, nil
}

func prepareCloudEventData(config *Config) Output {

	invokeCloudBeesWorkflow := &InvokeCloudBeesWorkflow{
		ComponentId:      config.ComponentId,
		BranchName:       config.BranchName,
		WorkflowFileName: config.WorkflowFileName,
		WorkflowInputs:   config.WorkflowInputs,
	}

	providerInfo := &ProviderInfo{
		RunId:      config.GhaRunId,
		RunAttempt: config.GhaRunAttempt,
		RunNumber:  config.GhaRunNumber,
		JobName:    config.GhaJobName,
		Provider:   GithubProvider,
	}
	output := Output{
		InvokeWorkflow: *invokeCloudBeesWorkflow,
		ProviderInfo:   *providerInfo,
	}
	fmt.Println("Output set data")
	fmt.Println(PrettyPrint(output))
	return output
}
func sendCloudEvent(cloudEvent cloudevents.Event, config *Config) error {
	eventJSON, err := json.Marshal(cloudEvent)
	if err != nil {
		return fmt.Errorf("error encoding CloudEvent JSON %s", err)
	}
	req, _ := http.NewRequest(PostMethod, getCloudbeesFullUrl(config), bytes.NewBuffer(eventJSON))
	fmt.Println(PrettyPrint(cloudEvent))
	// For Local Testing
	//req, _ := http.NewRequest(PostMethod, "http://localhost:8080/events", bytes.NewBuffer(eventJSON))

	req.Header.Set(ContentTypeHeaderKey, ContentTypeCloudEventsJson)
	req.Header.Set(AuthorizationHeaderKey, Bearer+config.CloudBeesApiToken)
	client := &http.Client{}
	resp, err := client.Do(req) // Fire and forget
	fmt.Println("Error from service:", err)
	fmt.Println("Response from service:", resp)
	fmt.Println(PrettyPrint(resp))

	if err != nil {
		return fmt.Errorf("error sending CloudEvent to platform %s", err)
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Println("received non-200 status code: %d, response: %s", resp.StatusCode, string(body))
		return fmt.Errorf("received non-200 status code: %d, response: %s", resp.StatusCode, string(body))
		// handle known error cases
		// switch resp.StatusCode {
		// // case http.StatusUnauthorized, http.StatusForbidden:
		// // 	return nil, constants.NoAccessErr
		// // case http.StatusNotFound:
		// // 	return nil, constants.NotFoundErr
		// // }
		//return fmt.Errorf("error sending CloudEvent to platform %s", resp.Status)
	}
	
	// Read the response body (so we can print it)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
	} else {
		// Print the response body
		fmt.Println("Response body:", string(body))
	}
	
	
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Error closing response body:", err)
		}
	}(resp.Body)

	fmt.Println("CloudEvent sent successfully!")
	return nil
}

// PrettyPrint converts the input to JSON string
func PrettyPrint(in interface{}) string {
	data, err := json.MarshalIndent(in, "", "  ")
	if err != nil {
		fmt.Println("error marshalling response", err)
	}
	return string(data)
}

func getCurrentBranchFromRef() (string, error) {
	githubRef := os.Getenv(GitHubRef)
	if githubRef == "" {
		return "", fmt.Errorf("GITHUB_REF environment variable is not set")
	}

	if strings.HasPrefix(githubRef, "refs/heads/") {
		return strings.TrimPrefix(githubRef, "refs/heads/"), nil
	}

	return "", fmt.Errorf("GITHUB_REF does not point to a branch, found: %s", githubRef)
}
