//go:build wireinject
// +build wireinject

package internal

import (
	"github.com/google/wire"
	"github.com/lunabrain-ai/lunapipe/internal/config"
	"github.com/urfave/cli/v2"
)

func Wire() (*cli.App, error) {
	panic(wire.Build(
		NewCLI,
		QAClientProviderSet,
		config.NewConfigProvider,
		config.NewOpenAIConfig,
		config.NewLogConfig,
	))
}
