package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/go-playground/validator/v10"
	"github.com/heka-ai/benchmark-api/internal/log"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

type APIConfig struct {
	config *Config
}

var logger = log.GetLogger("config")

var ConfigFX = fx.Module("config",
	fx.Provide(NewAPIConfig),
	fx.Invoke(func(c *APIConfig) {}),
)

func NewAPIConfig() *APIConfig {
	config := Init()
	apiConfig := &APIConfig{
		config: config,
	}

	apiConfig.WatchConfig()

	return apiConfig
}

func (c *APIConfig) GetConfig() *Config {
	return c.config
}

// Init the config and validate it
func Init() *Config {
	InitFlags()

	config := ReadConfig()

	return config
}

func InitFlags() {
	flag.String("config", "bench.toml", "Path to the config file")
	flag.Parse()
}

func ReadConfig() *Config {
	config := &Config{}
	filename := flag.Lookup("config").Value.String()

	viper.SetConfigName(filename)
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")

	viper.AutomaticEnv()
	viper.WatchConfig()

	err := viper.ReadInConfig()
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to read config file")
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to unmarshal config")
	}

	validate := validator.New()
	err = validate.Struct(config)
	if err != nil {
		// TODO: improve errors handling
		for _, err := range err.(validator.ValidationErrors) {
			logger.Error().Msgf("Failed to validate config, expected %s, got %s (tags: %s)", err.Field(), err.Value(), err.Tag())
			fmt.Println(err)
		}
		os.Exit(1)
	}

	logger.Info().Interface("config", config).Msgf("Config validated successfully")

	return config
}

func (c *APIConfig) WatchConfig() {
	viper.WatchConfig()

	viper.OnConfigChange(func(e fsnotify.Event) {
		newConfig := ReadConfig()

		c.config = newConfig

		logger.Info().Interface("config", c.config).Msgf("Config reloaded successfully")
	})
}
