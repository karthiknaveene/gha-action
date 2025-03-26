package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"gha-action/internal"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
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

	workflowInputs := os.Getenv(internal.WorkflowInputs)
	fmt.Println(workflowInputs)
	if workflowInputs != "" {
		err := json.Unmarshal([]byte(workflowInputs), &cfg.WorkflowInputs)
		if err != nil {
			fmt.Println("Error unmarshalling workflow inputs:", err)
		}
	} else {
		cfg.WorkflowInputs = make(map[string]string)
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
