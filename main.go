package main

import (
	"github.com/lunabrain-ai/lunabrain/pkg/store/cache"
	"github.com/lunabrain-ai/lunapipe/internal"
	"github.com/rs/zerolog/log"
	"os"
)

func main() {
	cacheConfig := cache.Config{
		Name: ".lunapipe",
	}

	app, err := internal.Wire(cacheConfig)
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
