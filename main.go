package main

import (
	"github.com/lunabrain-ai/aicli/internal"
	"github.com/rs/zerolog/log"
	"os"
)

func main() {
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
