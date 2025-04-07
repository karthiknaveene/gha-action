package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
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
		fmt.Printf("error preparing CloudEvent %v", err)
		return nil
	}
	err = sendCloudEvent(cloudEvent, config)
	if err != nil {
		fmt.Println(err)
		fmt.Printf("error sending CloudEvent %v", err)
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
func prepareCloudEvent(config *Config, cloudEventData CloudEventData) (cloudevents.Event, error) {
	cloudEvent := cloudevents.NewEvent()
	cloudEvent.SetID(uuid.NewString())
	cloudEvent.SetSubject(getSubject(config))
	cloudEvent.SetType(InvokeCloudBeesWorkflowEvent)
	cloudEvent.SetSource(getSource(config))
	cloudEvent.SetSpecVersion(SpecVersion)
	cloudEvent.SetTime(time.Now())
	err := cloudEvent.SetData(ContentTypeJson, cloudEventData)
	fmt.Println("CloudEvent set data")
	fmt.Println(PrettyPrint(cloudEvent))
	if err != nil {
		return cloudevents.Event{}, fmt.Errorf("failed to set data: %v", err)
	}
	return cloudEvent, nil
}

func prepareCloudEventData(config *Config) CloudEventData {

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
	cloudEventData := CloudEventData{
		InvokeWorkflow: *invokeCloudBeesWorkflow,
		ProviderInfo:   *providerInfo,
	}
	fmt.Println("Output set data")
	fmt.Println(PrettyPrint(cloudEventData))
	return cloudEventData
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

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
		}
		fmt.Println("Not successful response body - error code:", string(body))
		var errorResponse ErrorResponse
		if err := json.Unmarshal(body, &errorResponse); err != nil {
			return errors.New(string(body) + ".Please provide a valid cloudbees api url")
			//fmt.Println("Error unmarshaling response body:", err)
		}
		if errorResponse.Message == "" {
			return errors.New("Please provide a valid cloudbees api token")
		}
		return errors.New(errorResponse.Message)
	}

	// If status code is OK, print the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading successful response body: %s", err)
	}
	fmt.Println("Successful response body:", string(body))

	// Define the response structure based on the JSON format

	// Unmarshal the JSON into the struct
	var successResponse SuccessResponse
	if err := json.Unmarshal(body, &successResponse); err != nil {
		return fmt.Errorf("error unmarshaling response body: %s", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Error closing response body:", err)
		}
	}(resp.Body)

	//fmt.Printf(`::set-output name=cbp_run_url::%s`, successResponse.EventOutput.InvokeWorkflowOutput.RunUrl)
	// Output the runUrl to GITHUB_OUTPUT file for GitHub Actions
	runUrl := successResponse.EventOutput.InvokeWorkflowOutput.RunUrl
	err = writeGitHubOutput(runUrl)
	if err != nil {
	}

	fmt.Println("CloudEvent sent successfully!")
	if successResponse.ErrorMessage != "" {
		fmt.Printf("Error while invoking CloudBees workflow: %v", successResponse.ErrorMessage)
		//return fmt.Errorf("Error while invoking CloudBees workflow: %v", successResponse.ErrorMessage)
	}
	fmt.Printf("error %v", err)
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

	return "", fmt.Errorf("Please specify the branch of the CloudBees workflow, as the current GitHub workflow is not triggered by a branch.")
}

func writeGitHubOutput(runUrl string) error {
	// Open the GITHUB_OUTPUT file to append the output
	outputFile, err := os.OpenFile(os.Getenv("GITHUB_OUTPUT"), os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		fmt.Printf("Error opening GITHUB_OUTPUT file: %v", err)
		return nil
	}
	defer outputFile.Close()

	// Write the output to the GITHUB_OUTPUT file in the format expected by GitHub Actions
	_, err = fmt.Fprintf(outputFile, "cbp_run_url=%s\n", runUrl)
	if err != nil {
		fmt.Printf("Error writing to GITHUB_OUTPUT: %v", err)
	}
	return nil
}
