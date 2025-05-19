package main

import (
	"flag"
	"time"

	"github.com/getsentry/sentry-go"
	log "github.com/heka-ai/benchmark-cli/internal/logs"
)

var logger = log.GetLogger("cli")

func main() {
	disableTelemetry := flag.Bool("disable-telemetry", false, "Disable telemetry")
	flag.Parse()

	if !*disableTelemetry {
		err := sentry.Init(sentry.ClientOptions{
			Dsn: "https://0bf0fc25cd64524694b6dc78ada647e2@sentry.sia.partners/71",
		})
		if err != nil {
			logger.Error().Err(err).Msg("sentry.Init")
		}
		defer sentry.Flush(2 * time.Second)
	}

	rootCmd := RootCmd()
	rootCmd.Execute()
}
