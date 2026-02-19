package main

import (
	"github.com/spf13/cobra"

	cloud_generator "github.com/heka-ai/benchmark-cli/internal/cloud/generator"
	huggingface "github.com/heka-ai/benchmark-cli/internal/huggingface"
	"github.com/heka-ai/benchmark-cli/pkg/config"
)

// validate the cloud credentials, make sure the credentials are correct
func CredsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "creds",
		Short: "Validate the cloud credentials",
		Long:  `Validate the cloud credentials by running dry runs of the command we will use`,
		Run: func(cmd *cobra.Command, args []string) {
			excludedSections, err := cmd.Flags().GetStringArray("exclude-section")
			if err != nil {
				logger.Error().Err(err).Msg("Error getting excluded sections")
				return
			}
			validate(excludedSections)
		},
	}

	cmd.Flags().StringArray("exclude-section", []string{}, "Exclude a section from validation (can be used multiple times). Examples: --exclude-section benchmark --exclude-section vllm")

	return cmd
}

// Check if the credentials are valid (cloud and huggingface)
func validate(excludedSections []string) {
	logger.Info().Msg("Validating credentials")
	if len(excludedSections) > 0 {
		logger.Info().Strs("excluded_sections", excludedSections).Msg("Excluding sections from validation")
	}
	config.InitWithExclusions(excludedSections)
	c := config.GetConfig()

	cloud := cloud_generator.NewCloud(&c)
	if err := cloud.ValidateCredentials(); err != nil {
		logger.Fatal().Err(err).Msg("Cannot perform EC2 RunInstances dry-run")
	}

	if err := huggingface.ValidateHFCredentials(&c); err != nil {
		logger.Fatal().Err(err).Msg("HuggingFace credentials validation failed")
	}

	logger.Info().Msg("Credentials validated")
}
