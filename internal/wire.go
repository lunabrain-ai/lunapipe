//go:build wireinject
// +build wireinject

package internal

import (
	"github.com/google/wire"
	"github.com/lunabrain-ai/lunabrain/pkg/store/cache"
	"github.com/lunabrain-ai/lunapipe/internal/cli"
	"github.com/lunabrain-ai/lunapipe/internal/config"
	urfavcli "github.com/urfave/cli/v2"
)

func Wire(cacheConfig cache.Config) (*urfavcli.App, error) {
	panic(wire.Build(
		cli.New,
		config.ProviderSet,
	))
}
