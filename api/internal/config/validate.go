package apiConfig

import (
	"flag"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/fx"

	config "github.com/heka-ai/benchmark-cli/pkg/config"

	"github.com/heka-ai/benchmark-api/internal/log"
)

type APIConfig struct {
	config *config.Config
}

var logger = log.GetLogger("config")

var ConfigFX = fx.Module("config",
	fx.Provide(NewAPIConfig),
	fx.Invoke(func(c *APIConfig) {}),
)

func NewAPIConfig() *APIConfig {
	cfg := Init()
	apiCfg := &APIConfig{config: cfg}
	apiCfg.WatchConfig()
	return apiCfg
}

func (c *APIConfig) GetConfig() *config.Config {
	return c.config
}

// Init the config and validate it
func Init() *config.Config {
	InitFlags()
	return ReadConfig()
}

func InitFlags() {
	flag.String("config", "bench.toml", "Path to the config file")
	flag.Parse()
}

func ReadConfig() *config.Config {
	// Forward --config flag to the shared CLI config loader
	if f := flag.Lookup("config"); f != nil {
		filename := f.Value.String()
		if filename != "" {
			viper.Set("config", filename)
		}
	}

	// Deploy/infrastructure API startup should not require benchmark section.
	config.InitInfra()
	c := config.GetConfig()
	logger.Info().Interface("config", c).Msgf("Config validated successfully")
	return &c
}

func (c *APIConfig) WatchConfig() {
	viper.WatchConfig()

	viper.OnConfigChange(func(e fsnotify.Event) {
		newConfig := ReadConfig()

		c.config = newConfig

		logger.Info().Interface("config", c.config).Msgf("Config reloaded successfully")
	})
}
