package aws

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/spf13/viper"

	"github.com/heka-ai/benchmark-cli/internal/constants"
)

func (c *AWSClient) Create() error {
	filename := viper.GetString("config")
	if filename == "" {
		filename = "bench.toml"
	}
	configString, err := os.ReadFile(filename)

	if err != nil {
		logger.Error().Err(err).Msg("Error while reading the config file")
		return err
	}

	// set the basic benchmark config env var
	userData := fmt.Sprintf(`#!/bin/bash
echo "API_KEY=%s" > /home/ubuntu/.bashrc
echo "HF_TOKEN=%s" >> /home/ubuntu/.bashrc
cat > /home/ubuntu/config.toml << 'BENCH_EOF'
%s
BENCH_EOF
# Also export HF_TOKEN for the api.service (EnvironmentFile)
printf 'HF_TOKEN=%s\n' | sudo tee /etc/default/benchmark-api >/dev/null
sudo systemctl daemon-reload
sudo systemctl restart api
		`, c.config.APIKey, c.config.BenchmarkConfig.Token, configString, c.config.BenchmarkConfig.Token)

	// Optional networking from config
	subnetID := viper.GetString("aws.subnet_id")
	securityGroupID := viper.GetString("aws.security_group_id")

	// Helper to create tags (including defaults)
	buildTags := func(machineType string) []types.Tag {
		return []types.Tag{
			{Key: aws.String(constants.BenchInstanceLabelKey), Value: aws.String(machineType)},
			{Key: aws.String("managed-by"), Value: aws.String("benchmark-cli")},
			{Key: aws.String(constants.BenchmarkIDTag), Value: aws.String(c.config.BenchmarkID)},
			{Key: aws.String("Name"), Value: aws.String(fmt.Sprintf("benchmark-%s", c.config.BenchmarkID))},
		}
	}

	launchWithNetwork := func(instanceType, ami, machineType string) error {
		input := &ec2.RunInstancesInput{
			InstanceType: types.InstanceType(instanceType),
			ImageId:      aws.String(ami),
			MinCount:     aws.Int32(1),
			MaxCount:     aws.Int32(1),
			UserData:     aws.String(base64.StdEncoding.EncodeToString([]byte(userData))),
			TagSpecifications: []types.TagSpecification{{
				ResourceType: types.ResourceTypeInstance,
				Tags:         buildTags(machineType),
			}},
			BlockDeviceMappings: []types.BlockDeviceMapping{{
				DeviceName: aws.String("/dev/sda1"),
				Ebs:        &types.EbsBlockDevice{VolumeSize: aws.Int32(200)},
			}},
		}
		// Ensure a public IP is associated when explicit network is used
		iface := types.InstanceNetworkInterfaceSpecification{
			DeviceIndex:              aws.Int32(0),
			AssociatePublicIpAddress: aws.Bool(true),
			DeleteOnTermination:      aws.Bool(true),
		}
		if subnetID != "" {
			iface.SubnetId = aws.String(subnetID)
		}
		if securityGroupID != "" {
			iface.Groups = []string{securityGroupID}
		}
		input.NetworkInterfaces = []types.InstanceNetworkInterfaceSpecification{iface}

		_, err := c.svc.RunInstances(context.TODO(), input)
		if err != nil {
			return err
		}
		logger.Info().Msg("Instance created")
		return nil
	}

	useExplicitNetwork := subnetID != "" || securityGroupID != ""

	if useExplicitNetwork {
		// GPU/LLM instance
		if err = launchWithNetwork(c.config.AWSConfig.GPUInstanceType, c.config.AWSConfig.GPU_AMI, constants.LLMInstanceLabelValue); err != nil {
			logger.Error().Err(err).Msg("Error while creating the GPU instance")
			return err
		}
		// CPU/Bench instance
		if err = launchWithNetwork(c.config.AWSConfig.CPUInstanceType, c.config.AWSConfig.CPU_AMI, constants.BenchInstanceLabelValue); err != nil {
			logger.Error().Err(err).Msg("Error while creating the CPU instance")
			return err
		}
		return nil
	}

	// Fallback: default VPC and SG via existing helper
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
