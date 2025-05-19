package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	log "github.com/heka-ai/benchmark-cli/internal/logs"
	"github.com/spf13/viper"
)

var logger = log.GetLogger("config")
var config Config

// Init the config and validate it
func Init() {
	InitFlags()
	InitConfig()
}

func InitFlags() {
	flag.String("config", "bench.toml", "Path to the config file")
	flag.Parse()
}

func GetConfig() Config {
	return config
}

func InitConfig() {
	localConfig := Config{}

	filename := flag.Lookup("config").Value.String()

	viper.SetConfigName(filename)
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to read config file")
	}

	err = viper.Unmarshal(&localConfig)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to unmarshal config")
	}

	validate := validator.New()
	err = validate.Struct(localConfig)
	if err != nil {
		// TODO: improve errors handling
		for _, err := range err.(validator.ValidationErrors) {
			logger.Error().Msgf("Failed to validate config, expected %s, got %s (tags: %s)", err.Field(), err.Value(), err.Tag())
			fmt.Println(err)
		}
		os.Exit(1)
	}

	logger.Info().Msgf("Config validated successfully")

	config = localConfig
}
