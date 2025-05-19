package aws

import (
	"errors"
	"time"
)

func (c *AWSClient) Deploy() error {
	instances, err := c.GetBenchmarkInstances()
	if err != nil {
		logger.Error().Err(err).Msg("Cannot find instances link to the benchmark (have you run 'create' yet?)")
	}

	gpuInstanceIP := ""

	for _, instance := range instances {
		for _, tag := range instance.Tags {
			if *tag.Key == "benchmark-machine-type" && *tag.Value == "GPU" {
				gpuInstanceIP = *instance.PublicIpAddress
			}
		}

		if instance.PublicIpAddress == nil {
			logger.Fatal().Msg("Instance has no public IP address")
		}

		err := c.cli.Deploy(*instance.PublicIpAddress, c.config.InferenceEngine)
		logger.Info().Str("ip", *instance.PublicIpAddress).Msg("Deployment started")

		if err != nil {
			logger.Error().Err(err).Msg("Cannot deploy the model on the instance")
		}
	}

	maxRetries := 36
	for i := 0; i < maxRetries; i++ {
		ready, err := c.cli.ModelStatus(gpuInstanceIP)
		if err != nil {
			logger.Error().Err(err).Msg("Cannot get the model status")
		}
		if ready {
			logger.Info().Msg("Model is ready")
			return nil
		}
		time.Sleep(10 * time.Second)
	}

	return errors.New("model is not ready after 360 seconds")
}
