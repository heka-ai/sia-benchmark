package scaleway

import (
	"github.com/heka-ai/benchmark-cli/internal/bench"
	"github.com/heka-ai/benchmark-cli/internal/cloud"
	"github.com/heka-ai/benchmark-cli/internal/constants"
	log "github.com/heka-ai/benchmark-cli/internal/logs"
	"github.com/heka-ai/benchmark-cli/pkg/config"
	iam "github.com/scaleway/scaleway-sdk-go/api/iam/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/api/vpc/v2"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

var logger = log.GetLogger("scaleway")

type ScalewayClient struct {
	cloud.Cloud

	cli    *bench.Client
	config *config.Config

	iam      *iam.API
	vpc      *vpc.API
	instance *instance.API

	wasInit bool
}

func NewClient(config *config.Config) *ScalewayClient {
	return &ScalewayClient{
		cli:     bench.NewClient(config.APIKey),
		config:  config,
		wasInit: false,
	}
}

func (c *ScalewayClient) Init() cloud.Cloud {
	logger.Info().Msg("Initializing Scaleway client...")

	var opts []scw.ClientOption

	opts = append(opts, scw.WithDefaultOrganizationID(c.config.ScalewayConfig.OrganizationID))
	logger.Debug().Msg("Added Option WithDefaultOrganizationID... " + c.config.ScalewayConfig.OrganizationID)

	opts = append(opts, scw.WithDefaultProjectID(c.config.ScalewayConfig.ProjectID))
	logger.Debug().Msg("Added Option WithDefaultProjectID... " + c.config.ScalewayConfig.ProjectID)

	opts = append(opts, scw.WithDefaultRegion(scw.Region(c.config.ScalewayConfig.Region)))
	logger.Debug().Msg("Added Option WithDefaultRegion... " + c.config.ScalewayConfig.Region)

	if c.config.ScalewayConfig.ProfileName != "" {
		p, _ := scw.MustLoadConfig().GetProfile(c.config.ScalewayConfig.ProfileName)
		opts = append(opts, scw.WithProfile(p))
		logger.Debug().Msg("Added Option WithProfile... " + c.config.ScalewayConfig.ProfileName)
	} else {
		opts = append(opts, scw.WithAuth(c.config.ScalewayConfig.ScalewayAccessKey, c.config.ScalewayConfig.ScalewaySecretKey))
		logger.Debug().Msg("Added Option WithAuth... [SECRET KEYS NOT LOGGED]")
	}

	opts = append(opts, scw.WithEnv())
	logger.Debug().Msg("Added Option WithEnv... ")

	// Create a Scaleway client
	client, err := scw.NewClient(opts...)
	if err != nil {
		panic(err)
	}
	// Create SDK objects for Scaleway Instance product
	c.iam = iam.NewAPI(client)
	c.vpc = vpc.NewAPI(client)
	c.instance = instance.NewAPI(client)

	c.wasInit = true

	return c
}

func (c *ScalewayClient) GetAllBenchmarkServers() ([]*instance.Server, error) {
	return c.GetBenchmarkServers(nil)
}

func (c *ScalewayClient) GetBenchmarkServers(state *instance.ServerState) ([]*instance.Server, error) {
	listServersOutput, err := c.instance.ListServers(&instance.ListServersRequest{
		Tags:  []string{constants.BenchIDTag + "/" + c.config.BenchID},
		State: state,
	})

	if err != nil {
		return nil, err
	}

	return listServersOutput.Servers, nil
}

func (c *ScalewayClient) GetAllBenchmarkVolumes() ([]*instance.Volume, error) {
	listVolumesResponse, err := c.instance.ListVolumes(&instance.ListVolumesRequest{
		Tags: []string{constants.BenchIDTag + "/" + c.config.BenchID},
	})

	if err != nil {
		return nil, err
	}

	return listVolumesResponse.Volumes, nil
}

func (c *ScalewayClient) createServer() error {

	logger.Debug().Msgf("Creating Server... ")
	createServerResponse, err := c.instance.CreateServer(&instance.CreateServerRequest{
		CommercialType: "DEV1-S",
		Image:          scw.StringPtr("ubuntu_noble"),
		Tags:           []string{constants.BenchIDTag + "/" + c.config.BenchID},
		Project:        scw.StringPtr(c.config.ScalewayConfig.ProjectID),
	})
	if err != nil {
		return err
	}
	logger.Info().Msgf("Created Server %s (%s)... ", createServerResponse.Server.ID, createServerResponse.Server.Name)

	for _, volume := range createServerResponse.Server.Volumes {
		logger.Debug().Msgf("Updating Volume %s tags... ", volume.ID)
		_, err := c.instance.UpdateVolume(&instance.UpdateVolumeRequest{
			VolumeID: volume.ID,
			Tags:     &[]string{constants.BenchIDTag + "/" + c.config.BenchID},
		})
		if err != nil {
			return err
		}
		logger.Debug().Msgf("Updated Volume %s tags... ", volume.ID)
	}

	logger.Debug().Msgf("Using Action %s on Server %s... ", instance.ServerActionPoweron, createServerResponse.Server.ID)
	err = c.instance.ServerActionAndWait(&instance.ServerActionAndWaitRequest{
		ServerID: createServerResponse.Server.ID,
		Action:   instance.ServerActionPoweron,
	})
	if err != nil {
		return err
	}
	logger.Info().Msgf("Completed Action %s on Server %s (%s)... ", instance.ServerActionPoweron, createServerResponse.Server.ID, createServerResponse.Server.Name)

	return nil
}

func (c *ScalewayClient) deleteServer(server instance.Server) error {

	if server.State == instance.ServerStateRunning {
		logger.Debug().Msgf("Using Action %s on Server %s... ", instance.ServerActionPoweroff, server.ID)
		err := c.instance.ServerActionAndWait(&instance.ServerActionAndWaitRequest{
			ServerID: server.ID,
			Action:   instance.ServerActionPoweroff,
		})
		if err != nil {
			return err
		}
		logger.Debug().Msgf("Completed Action %s on Server %s... ", instance.ServerActionPoweroff, server.ID)
	}

	logger.Debug().Msgf("Deleting Server %s... ", server.ID)
	err := c.instance.DeleteServer(&instance.DeleteServerRequest{
		ServerID: server.ID,
	})
	if err != nil {
		return err
	}
	logger.Info().Msgf("Deleted Server %s... ", server.ID)

	return nil
}

func (c *ScalewayClient) deleteVolume(volume instance.Volume) error {

	logger.Debug().Msgf("Deleting Volume %s... ", volume.ID)
	err := c.instance.DeleteVolume(&instance.DeleteVolumeRequest{
		VolumeID: volume.ID,
	})
	if err != nil {
		return err
	}
	logger.Info().Msgf("Deleted Volume %s... ", volume.ID)

	return nil
}
