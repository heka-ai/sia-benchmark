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
			if tag.Key != nil && tag.Value != nil && *tag.Key == constants.BenchInstanceLabelKey && *tag.Value == constants.LLMInstanceLabelValue {
				if instance.PublicDnsName != nil && *instance.PublicDnsName != "" {
					return *instance.PublicDnsName, nil
				}
				if instance.PublicIpAddress != nil && *instance.PublicIpAddress != "" {
					return *instance.PublicIpAddress, nil
				}
				return "", errors.New("LLM instance has no public IP/DNS yet; try again shortly")
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
			if tag.Key != nil && tag.Value != nil && *tag.Key == constants.BenchInstanceLabelKey && *tag.Value == constants.BenchInstanceLabelValue {
				if instance.PublicDnsName != nil && *instance.PublicDnsName != "" {
					return *instance.PublicDnsName, nil
				}
				if instance.PublicIpAddress != nil && *instance.PublicIpAddress != "" {
					return *instance.PublicIpAddress, nil
				}
				return "", errors.New("CPU instance has no public IP/DNS yet; try again shortly")
			}
		}
	}
	return "", errors.New("no CPU instance found")
}
