package cloud_generator

import (
	"os"

	"github.com/heka-ai/benchmark-cli/internal/cloud"
	"github.com/heka-ai/benchmark-cli/internal/cloud/aws"
	"github.com/heka-ai/benchmark-cli/internal/config"
	log "github.com/heka-ai/benchmark-cli/internal/logs"
)

var logger = log.GetLogger("cloud_generator")

// NewCloud creates a new cloud client based on the provider
func NewCloud(config *config.Config) cloud.Cloud {
	switch config.Provider {
	case "aws":
		awsClient := aws.NewClient(config)
		return awsClient.Init()
	default:
		logger.Fatal().Msgf("Unsupported cloud provider: %s", config.Provider)
		os.Exit(1)
	}

	return nil
}
