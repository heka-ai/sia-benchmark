package main

import (
	cloud_generator "github.com/heka-ai/benchmark-cli/internal/cloud/generator"
	"github.com/heka-ai/benchmark-cli/pkg/config"
	"github.com/spf13/cobra"
)

// Destroy the instance created for the benchmark
// Remove all the data stored for the benchmark
func DestroyCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "destroy",
		Short: "Destroy the instances",
		Run: func(cmd *cobra.Command, args []string) {
			DestroyCmdExec()
		},
	}
}

func DestroyCmdExec() {
	logger.Info().Msg("Destroying the instance")
	config.Init()
	c := config.GetConfig()

	cloud := cloud_generator.NewCloud(&c)
	cloud.Destroy()

	logger.Info().Msg("Instances destroyed")
}
