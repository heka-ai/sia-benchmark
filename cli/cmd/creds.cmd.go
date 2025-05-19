package main

import (
	cloud_generator "github.com/heka-ai/benchmark-cli/internal/cloud/generator"
	"github.com/heka-ai/benchmark-cli/pkg/config"
	"github.com/spf13/cobra"
)

// validate the cloud credentials, make sure the credentials are correct
func CredsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "creds",
		Short: "Validate the cloud credentials",
		Long:  `Validate the cloud credentials by running dry runs of the command we will use`,
		Run: func(cmd *cobra.Command, args []string) {
			validate()
		},
	}
}

// Check if the credentials are valid (cloud and huggingface)
func validate() {
	logger.Info().Msg("Validating credentials")
	config.Init()
	c := config.GetConfig()

	cloud := cloud_generator.NewCloud(&c)
	cloud.ValidateCredentials()

	// todo: validate the huggingface cred

	logger.Info().Msg("Credentials validated")
}
