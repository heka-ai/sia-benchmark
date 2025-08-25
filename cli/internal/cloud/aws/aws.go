package aws

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"slices"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/smithy-go"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/heka-ai/benchmark-cli/internal/bench"
	"github.com/heka-ai/benchmark-cli/internal/cloud"
	"github.com/heka-ai/benchmark-cli/internal/constants"
	log "github.com/heka-ai/benchmark-cli/internal/logs"
	"github.com/heka-ai/benchmark-cli/pkg/config"
)

var logger = log.GetLogger("aws")

type AWSClient struct {
	cloud.Cloud

	cli     *bench.Client
	config  *config.Config
	svc     *ec2.Client
	wasInit bool
}

func NewClient(config *config.Config) *AWSClient {
	return &AWSClient{
		cli:     bench.NewClient(config.APIKey),
		config:  config,
		wasInit: false,
	}
}

func (c *AWSClient) Init() cloud.Cloud {
	logger.Info().Msg("Creating the AWS session using the credentials provided")
	var conf aws.Config
	var err error

	if c.config.AWSConfig.ProfileName != "" {
		logger.Debug().Str("profile", c.config.AWSConfig.ProfileName).Msg("Creating the EC2 client using this profile")
		conf, err = awsConfig.LoadDefaultConfig(context.TODO(),
			awsConfig.WithSharedConfigProfile(c.config.AWSConfig.ProfileName),
			awsConfig.WithRegion(c.config.AWSConfig.Region),
		)
	} else {
		conf, err = awsConfig.LoadDefaultConfig(context.TODO(),
			awsConfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(c.config.AWSConfig.AWSAccessKey, c.config.AWSConfig.AWSSecretKey, "")),
			awsConfig.WithRegion(c.config.AWSConfig.Region),
		)
	}

	if err != nil {
		logger.Error().Msgf("Failed to create AWS session: %v", err)
		os.Exit(1)
	}

	c.svc = ec2.NewFromConfig(conf)

	c.wasInit = true

	return c
}

func (c *AWSClient) validateCredentials() error {
	if !c.wasInit {
		return errors.New("client not initialized")
	}

	// Validate permissions by attempting a DryRun of RunInstances only
	err := c.createInstance(c.config.AWSConfig.CPUInstanceType, true, c.config.AWSConfig.GPU_AMI, []types.Tag{}, "")
	if err != nil {
		logger.Error().Err(err).Msg("Cannot perform EC2 RunInstances dry-run")
		return err
	}

	logger.Info().Msg("OK - Can create an instance")

	// this instance id do not exists
	err = c.deleteInstance("i-123456", true)
	if err != nil {
		logger.Error().Msgf("Cannot delete an EC2 instance")
		return err
	}

	logger.Info().Msg("OK - Can delete an instance")

	return nil
}

func isDryRunError(err error) bool {
	var apiErr smithy.APIError
	return errors.As(err, &apiErr) && apiErr.ErrorCode() == "DryRunOperation"
}

func (c *AWSClient) deleteInstance(instanceID string, dryRun bool) error {
	_, err := c.svc.TerminateInstances(context.TODO(), &ec2.TerminateInstancesInput{
		InstanceIds: []string{instanceID},
		DryRun:      aws.Bool(dryRun),
	})

	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) && dryRun && apiErr.ErrorCode() == "InvalidInstanceID.NotFound" {
			return nil
		}

		return err
	}

	return nil
}

func (c *AWSClient) createInstance(instanceType string, dryRun bool, ami string, tags []types.Tag, userData string) error {
	instanceName := fmt.Sprintf("benchmark-%s", c.config.BenchmarkID)

	defaultTags := []types.Tag{
		{
			Key:   aws.String("managed-by"),
			Value: aws.String("benchmark-cli"),
		},
		{
			Key:   aws.String("bench-id"),
			Value: aws.String(c.config.BenchmarkID),
		},
		{
			Key:   aws.String("Name"),
			Value: aws.String(instanceName),
		},
	}

	allTags := slices.Concat(tags, defaultTags)
	base64UserData := base64.StdEncoding.EncodeToString([]byte(userData))

	logger.Debug().Str("instance-type", instanceType).Str("ami", ami).Str("user-data", base64UserData).Interface("tags", allTags).Msg("Creating the instance (dry-run)")

	_, err := c.svc.RunInstances(context.TODO(), &ec2.RunInstancesInput{
		InstanceType: types.InstanceType(instanceType),
		ImageId:      aws.String(ami),
		MinCount:     aws.Int32(1),
		MaxCount:     aws.Int32(1),
		DryRun:       aws.Bool(dryRun),
		UserData:     aws.String(base64UserData),
		TagSpecifications: []types.TagSpecification{
			{
				ResourceType: types.ResourceTypeInstance,
				Tags:         allTags,
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

	if err != nil {
		if isDryRunError(err) {
			return nil
		}

		return err
	}

	return nil
}

// The following operations are unchanged for non-validation flows
func (c *AWSClient) GetBenchmarkInstances() ([]types.Instance, error) {
	describeInstanceOutput, err := c.svc.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{
		Filters: []types.Filter{
			{
				Name:   aws.String(fmt.Sprintf("tag:%s", constants.BenchmarkIDTag)),
				Values: []string{c.config.BenchmarkID},
			},
			{
				Name:   aws.String("instance-state-name"),
				Values: []string{"running"},
			},
		},
	})

	if err != nil {
		return nil, err
	}

	var allInstances []types.Instance

	for _, reservations := range describeInstanceOutput.Reservations {
		allInstances = append(allInstances, reservations.Instances...)
	}

	return allInstances, nil
}

func (c *AWSClient) CreateInstance(instanceType string, ami string, tags []types.Tag, userData string) error {
	if !c.wasInit {
		return errors.New("client not initialized")
	}

	err := c.createInstance(instanceType, false, ami, tags, userData)
	if err != nil {
		return err
	}

	logger.Info().Msg("Instance created")

	return nil
}

func (c *AWSClient) DeleteInstance() error {
	if !c.wasInit {
		return errors.New("client not initialized")
	}

	instances, err := c.GetBenchmarkInstances()

	if err != nil {
		return err
	}

	if len(instances) == 0 {
		return errors.New("no instance found for this benchmark")
	}

	for _, v := range instances {
		err = c.deleteInstance(*v.InstanceId, false)

		if err != nil {
			return err
		}

		logger.Info().Msg("Instance terminated")
	}

	return nil
}
