package aws

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/spf13/viper"

	"github.com/heka-ai/benchmark-cli/internal/constants"
)

func (c *AWSClient) Create(createLLM, createBench bool) error {
	filename := viper.GetString("config")
	if filename == "" {
		filename = "bench.toml"
	}
	configString, err := os.ReadFile(filename)

	if err != nil {
		logger.Error().Err(err).Msg("Error while reading the config file")
		return err
	}

	hfToken := ""
	if c.config.BenchmarkConfig != nil && strings.TrimSpace(c.config.BenchmarkConfig.Token) != "" {
		hfToken = strings.TrimSpace(c.config.BenchmarkConfig.Token)
	} else if envToken := strings.TrimSpace(os.Getenv("HF_TOKEN")); envToken != "" {
		hfToken = envToken
	} else {
		logger.Warn().Msg("No HF token found in config or HF_TOKEN env var; set it before deploy")
	}

	userData := fmt.Sprintf(`#!/bin/bash
set -euo pipefail
exec > /var/log/bench-bootstrap.log 2>&1

export DEBIAN_FRONTEND=noninteractive

echo "=== Writing config ==="
cat > /home/ubuntu/config.toml << 'BENCH_EOF'
%s
BENCH_EOF
chown ubuntu:ubuntu /home/ubuntu/config.toml

echo "=== Setting env vars ==="
cat > /home/ubuntu/.benchrc << 'ENVEOF'
API_KEY=%s
HF_TOKEN=%s
PATH=/home/ubuntu/.local/bin:/usr/local/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
ENVEOF
chown ubuntu:ubuntu /home/ubuntu/.benchrc
echo 'source /home/ubuntu/.benchrc' >> /home/ubuntu/.bashrc

echo "=== Installing system deps ==="
apt-get update -y
apt-get install -y git curl jq

echo "=== Installing Go ==="
curl -fsSL https://go.dev/dl/go1.24.7.linux-amd64.tar.gz | tar -C /usr/local -xz

echo "=== Cloning repo ==="
sudo -u ubuntu git clone https://github.com/heka-ai/sia-benchmark.git /home/ubuntu/sia-benchmark

echo "=== Building API binary ==="
sudo -u ubuntu bash -c 'export PATH=/usr/local/go/bin:$PATH && cd /home/ubuntu/sia-benchmark && go build -o /home/ubuntu/api ./api/cmd'

echo "=== Installing systemd service ==="
cat > /etc/systemd/system/api.service << 'SVCEOF'
[Unit]
Description=Benchmark API Server
After=network.target

[Service]
WorkingDirectory=/home/ubuntu
ExecStart=/home/ubuntu/api --config /home/ubuntu/config.toml
Restart=always
RestartSec=5
User=ubuntu
Group=ubuntu
EnvironmentFile=/home/ubuntu/.benchrc

[Install]
WantedBy=multi-user.target
SVCEOF

systemctl daemon-reload
systemctl enable api.service
systemctl start api.service

echo "=== Bootstrap complete ==="
		`, configString, c.config.GeneralConfig.APIKey, hfToken)

	// Optional networking from config
	subnetID := viper.GetString("aws.subnet_id")
	securityGroupID := viper.GetString("aws.security_group_id")

	// Helper to create tags (including defaults)
	buildTags := func(machineType string) []types.Tag {
		return []types.Tag{
			{Key: aws.String(constants.BenchInstanceLabelKey), Value: aws.String(machineType)},
			{Key: aws.String("managed-by"), Value: aws.String("benchmark-cli")},
			{Key: aws.String(constants.BenchmarkIDTag), Value: aws.String(c.config.GeneralConfig.BenchmarkID)},
			{Key: aws.String("Name"), Value: aws.String(fmt.Sprintf("benchmark-%s", c.config.GeneralConfig.BenchmarkID))},
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
		if createLLM {
			if err = launchWithNetwork(c.config.AWSConfig.GPUInstanceType, c.config.AWSConfig.AMI, constants.LLMInstanceLabelValue); err != nil {
				logger.Error().Err(err).Msg("Error while creating the GPU instance")
				return err
			}
		}
		if createBench {
			if err = launchWithNetwork(c.config.AWSConfig.CPUInstanceType, c.config.AWSConfig.AMI, constants.BenchInstanceLabelValue); err != nil {
				logger.Error().Err(err).Msg("Error while creating the CPU instance")
				return err
			}
		}
		return nil
	}

	if createLLM {
		err = c.CreateInstance(c.config.AWSConfig.GPUInstanceType, c.config.AWSConfig.AMI, []types.Tag{
			{Key: aws.String(constants.BenchInstanceLabelKey), Value: aws.String(constants.LLMInstanceLabelValue)},
		}, userData)
		if err != nil {
			logger.Error().Err(err).Msg("Error while creating the GPU instance")
			return err
		}
	}
	if createBench {
		err = c.CreateInstance(c.config.AWSConfig.CPUInstanceType, c.config.AWSConfig.AMI, []types.Tag{
			{Key: aws.String(constants.BenchInstanceLabelKey), Value: aws.String(constants.BenchInstanceLabelValue)},
		}, userData)
		if err != nil {
			logger.Error().Err(err).Msg("Error while creating the CPU instance")
			return err
		}
	}

	return nil
}
