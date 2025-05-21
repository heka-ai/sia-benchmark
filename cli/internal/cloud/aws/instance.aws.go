package aws

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/google/uuid"
	"github.com/heka-ai/benchmark-cli/internal/constants"
	"github.com/heka-ai/benchmark-cli/internal/ssh"
)

const BASE_AMI_ID = "ami-009604998d7aa26d4"
const AMI_INSTANCE_TEMPLATE_TAG = "ami-instance-template"

func (c *AWSClient) CreateTemplateInstance() error {
	keyPairName := fmt.Sprintf("benchmark-key-pair-%s", uuid.New().String())
	logger.Debug().Str("keyName", keyPairName).Msg("Generating a key pair")

	keyPair, err := c.svc.CreateKeyPair(context.TODO(), &ec2.CreateKeyPairInput{
		KeyName: aws.String(keyPairName),
	})

	if err != nil {
		return err
	}

	defer func() {
		logger.Debug().Str("keyName", keyPairName).Msg("Deleting the key pair")
		c.svc.DeleteKeyPair(context.TODO(), &ec2.DeleteKeyPairInput{
			KeyName: aws.String(keyPairName),
		})
		logger.Debug().Str("keyName", keyPairName).Msg("Key pair deleted from AWS")
	}()

	if keyPair.KeyMaterial == nil {
		return fmt.Errorf("failed to create key pair")
	}

	// write the key to a temp file
	keyFile, err := os.CreateTemp("/tmp", "tmp-key-ssh-benchmark-*.pem")

	if err != nil {
		return err
	}

	defer func() {
		logger.Debug().Str("keyName", keyPairName).Msg("Deleting the key pair file")
		os.Remove(keyFile.Name())
	}()

	// write the key to the file
	_, err = keyFile.WriteString(*keyPair.KeyMaterial)

	if err != nil {
		return err
	}

	logger.Debug().Str("keyName", keyPairName).Msg("Key pair generated")
	logger.Debug().Str("instanceType", c.config.AWSConfig.GPUInstanceType).Msg("Creating the template instance")

	instance, err := c.svc.RunInstances(context.TODO(), &ec2.RunInstancesInput{
		ImageId:      aws.String(BASE_AMI_ID),
		InstanceType: types.InstanceType(c.config.AWSConfig.GPUInstanceType),
		MinCount:     aws.Int32(1),
		MaxCount:     aws.Int32(1),
		KeyName:      aws.String(keyPairName),
		TagSpecifications: []types.TagSpecification{
			{
				ResourceType: types.ResourceTypeInstance,
				Tags:         []types.Tag{{Key: aws.String(constants.BenchInstanceLabelKey), Value: aws.String(AMI_INSTANCE_TEMPLATE_TAG)}},
			},
		},
		BlockDeviceMappings: []types.BlockDeviceMapping{
			{
				DeviceName: aws.String("/dev/sda1"),
				Ebs: &types.EbsBlockDevice{
					VolumeSize: aws.Int32(200),
				},
			},
		},
	})

	defer func() {
		logger.Debug().Str("instanceId", *instance.Instances[0].InstanceId).Msg("Terminating the instance")
		c.svc.TerminateInstances(context.TODO(), &ec2.TerminateInstancesInput{
			InstanceIds: []string{*instance.Instances[0].InstanceId},
		})
		logger.Debug().Str("instanceId", *instance.Instances[0].InstanceId).Msg("Instance terminated")
	}()

	if err != nil {
		return err
	}

	logger.Debug().Str("keyMaterial", *keyPair.KeyMaterial).Msg("Creating the SSH client")

	waiter := ec2.NewInstanceStatusOkWaiter(c.svc)

	logger.Info().Msg("Waiting for the instance to be running")
	err = waiter.Wait(context.TODO(), &ec2.DescribeInstanceStatusInput{
		InstanceIds: []string{*instance.Instances[0].InstanceId},
	}, 10*time.Minute)

	if err != nil {
		return err
	}

	logger.Info().Msg("Instance is running")

	// describe the instance to get the public ip address
	describeInstance, err := c.svc.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{
		InstanceIds: []string{*instance.Instances[0].InstanceId},
	})

	if err != nil {
		return err
	}

	logger.Debug().Str("publicIp", *describeInstance.Reservations[0].Instances[0].PublicIpAddress).Str("keyFile", keyFile.Name()).Msg("Creating the SSH client")

	sshClient := ssh.NewSSHClient(keyFile.Name(), *describeInstance.Reservations[0].Instances[0].PublicIpAddress, "ubuntu")

	output, err := sshClient.Run("ls /tmp/")

	if err != nil {
		logger.Error().Err(err).Msg("Failed to run command on the instance")
		return err
	}

	logger.Info().Msgf("Output: %s", output)

	return nil
}

func (c *AWSClient) CreateBenchInstance() error {
	logger.Info().Msg("Starting the creation of the CPU instance image")

	err := c.CreateTemplateInstance()

	return err
}

func (c *AWSClient) CreateLLMInstance() error {
	logger.Info().Msg("Starting the creation of the LLM instance image")

	err := c.CreateTemplateInstance()

	return err
}
