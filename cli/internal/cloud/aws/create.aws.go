package aws

import (
	"flag"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/heka-ai/benchmark-cli/internal/constants"
)

func (c *AWSClient) Create() error {
	filename := flag.Lookup("config").Value.String()
	configString, err := os.ReadFile(filename)

	if err != nil {
		logger.Error().Err(err).Msg("Error while reading the config file")
		return err
	}

	// set the basic benchmark config env var
	userData := fmt.Sprintf(`#!/bin/bash
echo "API_KEY=%s" > /home/ubuntu/.bashrc
TOML_FILE='
%s
'
touch /home/ubuntu/config.toml
echo "$TOML_FILE" > /home/ubuntu/config.toml
		`, c.config.APIKey, configString)

	err = c.CreateInstance(c.config.AWSConfig.GPUInstanceType, c.config.AWSConfig.GPU_AMI, []types.Tag{
		{
			Key:   aws.String(constants.BenchInstanceLabelKey),
			Value: aws.String(constants.LLMInstanceLabelValue),
		},
	}, userData)

	if err != nil {
		logger.Error().Err(err).Msg("Error while creating the GPU instance")
		return err
	}

	err = c.CreateInstance(c.config.AWSConfig.CPUInstanceType, c.config.AWSConfig.CPU_AMI, []types.Tag{
		{
			Key:   aws.String(constants.BenchInstanceLabelKey),
			Value: aws.String(constants.BenchInstanceLabelValue),
		},
	}, userData)

	if err != nil {
		// TODO: delete the GPU instance
		logger.Error().Err(err).Msg("Error while creating the CPU instance")
		return err
	}

	return nil
}
