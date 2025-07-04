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

func (c *ScalewayClient) deleteServer(serverID string) error {

	action := instance.ServerActionPoweroff

	logger.Debug().Msgf("Using Action %s on Server %s... ", action, serverID)
	err := c.instance.ServerActionAndWait(&instance.ServerActionAndWaitRequest{
		ServerID: serverID,
		Action:   action,
	})
	if err != nil {
		return err
	}
	logger.Debug().Msgf("Completed Action %s on Server %s... ", action, serverID)

	err = c.instance.DeleteServer(&instance.DeleteServerRequest{
		ServerID: serverID,
	})
	if err != nil {
		return err
	}
	logger.Debug().Msgf("Deleted Server %s... ", serverID)

	return nil
}

func (c *ScalewayClient) deleteVolume(volumeID string) error {

	err := c.instance.DeleteVolume(&instance.DeleteVolumeRequest{
		VolumeID: volumeID,
	})
	if err != nil {
		return err
	}
	logger.Debug().Msgf("Deleted Volume %s... ", volumeID)

	return nil
}
