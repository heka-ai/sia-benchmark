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

func (c *AWSClient) CreateTemplateInstance(installType string) error {
	keyPairName := fmt.Sprintf("benchmark-key-pair-%s", uuid.New().String())
	logger.Debug().Str("keyName", keyPairName).Msg("Generating a key pair")

	keyPair, err := c.svc.CreateKeyPair(context.TODO(), &ec2.CreateKeyPairInput{
		KeyName: aws.String(keyPairName),
	})

	if err != nil {
		return err
	}

	// defer func() {
	// 	logger.Debug().Str("keyName", keyPairName).Msg("Deleting the key pair")
	// 	c.svc.DeleteKeyPair(context.TODO(), &ec2.DeleteKeyPairInput{
	// 		KeyName: aws.String(keyPairName),
	// 	})
	// 	logger.Debug().Str("keyName", keyPairName).Msg("Key pair deleted from AWS")
	// }()

	if keyPair.KeyMaterial == nil {
		return fmt.Errorf("failed to create key pair")
	}

	// write the key to a temp file
	keyFile, err := os.CreateTemp("/tmp", "tmp-key-ssh-benchmark-*.pem")

	if err != nil {
		return err
	}

	// defer func() {
	// 	logger.Debug().Str("keyName", keyPairName).Msg("Deleting the key pair file")
	// 	os.Remove(keyFile.Name())
	// }()

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

	// defer func() {
	// 	logger.Debug().Str("instanceId", *instance.Instances[0].InstanceId).Msg("Terminating the instance")
	// 	c.svc.TerminateInstances(context.TODO(), &ec2.TerminateInstancesInput{
	// 		InstanceIds: []string{*instance.Instances[0].InstanceId},
	// 	})
	// 	logger.Debug().Str("instanceId", *instance.Instances[0].InstanceId).Msg("Instance terminated")
	// }()

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

	err = c.runSetup(sshClient, installType)

	if err != nil {
		return err
	}

	logger.Info().Msg("Setup complete, creating the AMI")

	amiName := fmt.Sprintf("benchmark-ami-%s-%s", installType, time.Now().Format("2006-01-02-15-04-05"))

	ami, err := c.svc.CreateImage(context.TODO(), &ec2.CreateImageInput{
		InstanceId:  instance.Instances[0].InstanceId,
		Name:        aws.String(amiName),
		Description: aws.String(fmt.Sprintf("AMI for the benchmark-cli %s instance", installType)),
	})

	if err != nil {
		return err
	}

	imageAvailableWaiter := ec2.NewImageAvailableWaiter(c.svc)

	err = imageAvailableWaiter.Wait(context.TODO(), &ec2.DescribeImagesInput{
		ImageIds: []string{*ami.ImageId},
	}, 30*time.Minute)

	if err != nil {
		return err
	}
	logger.Info().Str("amiId", *ami.ImageId).Msg("AMI created")

	return nil
}

func (c *AWSClient) runSetup(sshClient *ssh.SSHClient, installType string) error {
	// clone the benchmark-cli repo
	logger.Info().Msg("Installing the benchmark-cli repo")
	err := sshClient.Run("git clone https://github.com/heka-ai/sia-benchmark.git /home/ubuntu/benchmark-cli")

	if err != nil {
		return err
	}

	err = sshClient.Run("mv /home/ubuntu/benchmark-cli/instance-builder/aws/ec2 /home/ubuntu/")

	if err != nil {
		return err
	}

	sshClient.Run("ls /home/ubuntu")

	logger.Info().Msg("Installing the control API")
	err = sshClient.Run("sudo bash /home/ubuntu/ec2/api/install.sh")

	if err != nil {
		return err
	}

	if installType == "llm" {
		logger.Info().Msg("Installing the LLM dependencies")
		err = sshClient.Run("sudo bash /home/ubuntu/ec2/gpu/install.sh")
	} else {
		logger.Info().Msg("Installing the CPU dependencies")
		err = sshClient.Run("sudo bash /home/ubuntu/ec2/cpu/install.sh")
	}

	return err
}

func (c *AWSClient) CreateBenchInstance() error {
	logger.Info().Msg("Starting the creation of the CPU instance image")

	err := c.CreateTemplateInstance("cpu")

	return err
}

func (c *AWSClient) CreateLLMInstance() error {
	logger.Info().Msg("Starting the creation of the LLM instance image")

	err := c.CreateTemplateInstance("llm")

	return err
}
