package main

import (
	cloud_generator "github.com/heka-ai/benchmark-cli/internal/cloud/generator"
	"github.com/heka-ai/benchmark-cli/pkg/config"
	"github.com/spf13/cobra"
)

func InstanceBuildCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create-instance [instance-type]",
		Short: "Create a the instance image for the cloud provider",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			instanceType := args[0]
			buildImageInstance(instanceType)
		},
	}
}

func buildImageInstance(instanceType string) {
	config.Init()
	config := config.GetConfig()

	cloud := cloud_generator.NewCloud(&config)

	var err error

	if instanceType == "llm" {
		err = cloud.CreateLLMInstance()
	} else if instanceType == "bench" {
		err = cloud.CreateBenchInstance()
	} else {
		logger.Error().Str("instance-type", instanceType).Msg("Invalid instance type")
	}

	if err != nil {
		logger.Error().Err(err).Msg("Error occurred while creating the instance")
	}
}
