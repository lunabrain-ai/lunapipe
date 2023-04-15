package config

import (
	"github.com/google/wire"
	"github.com/lunabrain-ai/lunabrain/pkg/store/cache"
)

var ProviderSet = wire.NewSet(
	cache.NewLocalCache,
	wire.Bind(new(cache.Cache), new(*cache.LocalCache)),

	NewConfigProvider,
	NewOpenAIConfig,
	NewLogConfig,
	NewConfigurator,
	wire.Bind(new(Configurator), new(*LocalConfigurator)),
)
