package scw

import (
	"github.com/heka-ai/benchmark-cli/internal/bench"
	"github.com/heka-ai/benchmark-cli/internal/cloud"
	log "github.com/heka-ai/benchmark-cli/internal/logs"
	"github.com/heka-ai/benchmark-cli/pkg/config"
)

var logger = log.GetLogger("scw")

type SCWClient struct {
	cloud.Cloud

	cli     *bench.Client
	config  *config.Config
	wasInit bool
}

func NewClient(config *config.Config) *SCWClient {
	return &SCWClient{
		cli:     bench.NewClient(config.APIKey),
		config:  config,
		wasInit: false,
	}
}

func (c *SCWClient) Init() cloud.Cloud {
	logger.Info().Msg("Mocking Scaleway session")

	return c
}
