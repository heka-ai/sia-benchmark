package main

import (
	cloud_generator "github.com/heka-ai/benchmark-cli/internal/cloud/generator"
	"github.com/heka-ai/benchmark-cli/pkg/config"
	"github.com/spf13/cobra"
)

func InstanceCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create",
		Short: "Create the instance to run the benchmark",
		Run: func(cmd *cobra.Command, args []string) {
			create()
		},
	}
}

func create() {
	logger.Info().Msg("Creating the instances to run the benchmark")
	config.Init()
	c := config.GetConfig()

	cloud := cloud_generator.NewCloud(&c)
	cloud.Create()

	// todo: add a flag to wait for the instances to be ready

	logger.Info().Msg("Instances created")
}
