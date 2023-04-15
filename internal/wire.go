//go:build wireinject
// +build wireinject

package internal

import (
	"github.com/google/wire"
	"github.com/lunabrain-ai/lunabrain/pkg/store/cache"
	"github.com/lunabrain-ai/lunapipe/internal/config"
	"github.com/urfave/cli/v2"
)

func Wire(cacheConfig cache.Config) (*cli.App, error) {
	panic(wire.Build(
		NewCLI,
		QAClientProviderSet,
		config.ProviderSet,
	))
}
