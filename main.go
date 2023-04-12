package main

import (
	"github.com/UnnoTed/horizontal"
	"github.com/lunabrain-ai/aicli/internal"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"strings"
)

// TODO breadchris this should be a provided dependency
func setupLogging() {
	logLevel := zerolog.InfoLevel
	if strings.ToLower(os.Getenv("LOG_LEVEL")) == "debug" {
		logLevel = zerolog.DebugLevel
	}
	log.Logger = zerolog.New(horizontal.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger().Level(logLevel)
}

func main() {
	setupLogging()

	app, err := internal.Wire()
	if err != nil {
		log.Error().Msgf("%+v\n", err)
		return
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Error().Msgf("%+v\n", err)
		return
	}
}
