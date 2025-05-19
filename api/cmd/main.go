package main

import (
	injection "github.com/heka-ai/benchmark-api/internal/fx"
	"github.com/heka-ai/benchmark-api/internal/log"
)

var logger = log.GetLogger("main")

func main() {
	logger.Info().Msg("Starting the communication API")

	injector := injection.NewAppInjector()
	injector.Start()
}
