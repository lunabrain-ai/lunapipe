package config

import (
	"github.com/lunabrain-ai/lunabrain/pkg/store/cache"
	"github.com/lunabrain-ai/lunapipe/internal/log"
	"github.com/lunabrain-ai/lunapipe/internal/openai"
	"go.uber.org/config"
	"os"
	"path"
)

const (
	localConfigFile = ".lunapipe.yaml"
	homeConfigFile  = "config.yaml"
)

type BaseConfig struct {
	Cache  cache.Config  `yaml:"cache"`
	OpenAI openai.Config `yaml:"openai"`
	Log    log.Config    `yaml:"log"`
}

func NewDefaultConfig() BaseConfig {
	return BaseConfig{
		Cache: cache.Config{
			Name: ".lunapipe",
		},
		OpenAI: openai.NewDefaultConfig(),
		Log:    log.NewDefaultConfig(),
	}
}

func NewProvider(cache cache.Cache) (config.Provider, error) {
	opts := []config.YAMLOption{
		config.Permissive(),
		config.Expand(os.LookupEnv),
		config.Static(NewDefaultConfig()),
	}

	if f, err := os.Stat(localConfigFile); err == nil {
		opts = append(opts, config.File(path.Join(f.Name())))
	}

	homeFile, err := cache.GetFile(homeConfigFile)
	if err == nil {
		if _, err := os.Stat(homeFile); err == nil {
			opts = append(opts, config.File(homeFile))
		}
	}
	return config.NewYAML(opts...)
}
