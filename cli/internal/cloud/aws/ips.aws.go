package aws

import (
	"errors"

	"github.com/heka-ai/benchmark-cli/internal/constants"
)

func (c *AWSClient) GetLLMInstanceIP() (string, error) {
	instances, err := c.GetBenchmarkInstances()
	if err != nil {
		return "", err
	}

	for _, instance := range instances {
		for _, tag := range instance.Tags {
			if *tag.Key == constants.BenchInstanceLabelKey && *tag.Value == constants.LLMInstanceLabelValue {
				return *instance.PublicIpAddress, nil
			}
		}
	}

	return "", errors.New("no LLM instance found")
}

func (c *AWSClient) GetBenchInstanceIP() (string, error) {
	instances, err := c.GetBenchmarkInstances()
	if err != nil {
		return "", err
	}

	for _, instance := range instances {
		for _, tag := range instance.Tags {
			if *tag.Key == constants.BenchInstanceLabelKey && *tag.Value == constants.BenchInstanceLabelValue {
				return *instance.PublicIpAddress, nil
			}
		}
	}

	return "", errors.New("no CPU instance found")
}
