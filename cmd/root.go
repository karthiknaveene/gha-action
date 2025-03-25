package cmd

import (
	"context"
	"fmt"
	"gha-action/internal"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"strings"
)

var (
	cmd = &cobra.Command{
		Use:   "gha-run-cbp-workflow",
		Short: "Trigger CloudBees workflow from GitHub Actions workflow",
		Long:  "Trigger CloudBees workflow from GitHub Actions workflow",
		RunE:  run,
	}
	cfg internal.Config
)

func Execute() error {
	return cmd.Execute()
}

func init() {
	setDefaultValues(&cfg)
}

func setDefaultValues(cfg *internal.Config) {
	componentId := os.Getenv(internal.ComponentId)
	if componentId != "" {
		cfg.ComponentId = componentId
	} else {
		cfg.ComponentId = ""
	}

	branchName := os.Getenv(internal.BranchName)
	if branchName != "" {
		cfg.BranchName = branchName
	} else {
		currentBranch, _ := getCurrentBranchFromRef()
		cfg.BranchName = currentBranch
	}

	workflowFileName := os.Getenv(internal.WorkflowFileName)
	if workflowFileName != "" {
		cfg.WorkflowFileName = workflowFileName
	} else {
		cfg.WorkflowFileName = ""
	}

	workflowInputs := os.Getenv(internal.WorkflowInputs)
	if workflowInputs != "" {
		cfg.WorkflowInputs = workflowInputs
	} else {
		cfg.WorkflowInputs = ""
	}
}

func run(_ *cobra.Command, args []string) error {
	if len(args) > 0 {
		return fmt.Errorf("unknown arguments: %v", args)
	}
	newContext, cancel := context.WithCancel(context.Background())
	osChannel := make(chan os.Signal, 1)
	signal.Notify(osChannel, os.Interrupt)
	go func() {
		<-osChannel
		cancel()
	}()

	return cfg.Run(newContext)
}

func getCurrentBranchFromRef() (string, error) {
	githubRef := os.Getenv("GITHUB_REF")
	if githubRef == "" {
		return "", fmt.Errorf("GITHUB_REF environment variable is not set")
	}

	if strings.HasPrefix(githubRef, "refs/heads/") {
		return strings.TrimPrefix(githubRef, "refs/heads/"), nil
	}

	return "", fmt.Errorf("GITHUB_REF does not point to a branch, found: %s", githubRef)
}
