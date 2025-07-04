package scaleway

import "github.com/scaleway/scaleway-sdk-go/api/vpc/v2"

func (c *ScalewayClient) Create() error {
	logger.Info().Msg("Mocking Scaleway Create")

	found_vpc_list, _ := c.vpc.ListVPCs(&vpc.ListVPCsRequest{})
	for _, found_vpc := range found_vpc_list.Vpcs {
		logger.Info().Msg("Found VPC: " + found_vpc.Name + " (" + found_vpc.ID + ")")
	}

	return nil
}
