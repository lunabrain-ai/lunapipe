package config

import (
	"github.com/google/wire"
	"github.com/lunabrain-ai/lunabrain/pkg/store/cache"
	"github.com/lunabrain-ai/lunapipe/internal/log"
	"github.com/lunabrain-ai/lunapipe/internal/openai"
)

var ProviderSet = wire.NewSet(
	cache.NewLocalCache,
	wire.Bind(new(cache.Cache), new(*cache.LocalCache)),

	NewConfigurator,
	wire.Bind(new(Configurator), new(*LocalConfigurator)),

	openai.NewConfig,
	log.NewConfig,

	NewProvider,
)
