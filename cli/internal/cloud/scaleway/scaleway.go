package scaleway

import (
	"github.com/heka-ai/benchmark-cli/internal/bench"
	"github.com/heka-ai/benchmark-cli/internal/cloud"
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

	return c
}
