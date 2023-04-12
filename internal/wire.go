//go:build wireinject
// +build wireinject

package internal

import (
	"github.com/google/wire"
	"github.com/urfave/cli/v2"
)

func Wire() (*cli.App, error) {
	panic(wire.Build(
		NewCLI,
		QAClientProviderSet,
		NewConfigProvider,
		NewOpenAIConfig,
		NewLogConfig,
	))
}
