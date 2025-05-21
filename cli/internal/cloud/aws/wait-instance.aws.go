package aws

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

const (
	MAX_ITERATIONS = 100
)

func (c *AWSClient) WaitForInstances() error {
	instances, err := c.GetBenchmarkInstances()
	if err != nil {
		return err
	}

	instanceIds := make([]string, len(instances))
	for i, instance := range instances {
		instanceIds[i] = *instance.InstanceId
	}

	waiter := ec2.NewInstanceRunningWaiter(c.svc)

	err = waiter.Wait(context.TODO(), &ec2.DescribeInstancesInput{
		InstanceIds: instanceIds,
	}, 10*time.Minute)

	if err != nil {
		return err
	}

	return nil
}
