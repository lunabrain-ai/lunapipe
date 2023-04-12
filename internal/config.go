package internal

import (
	"github.com/rs/zerolog/log"
	"go.uber.org/config"
	"os"
	"path"
)

const configFile = ".aicli.yaml"

type Config struct {
	APIKey string `yaml:"api_key"`
}

type BaseConfig struct {
	OpenAI Config `yaml:"openai"`
}

func NewCLIConfig(provider config.Provider) (Config, error) {
	var c Config
	err := provider.Get("openai").Populate(&c)
	if err != nil {
		return Config{}, err
	}
	return c, nil
}

func newDefaultConfig() BaseConfig {
	return BaseConfig{
		OpenAI: Config{
			APIKey: "${OPENAI_API_KEY}",
		},
	}
}

func NewConfigProvider() (config.Provider, error) {
	opts := []config.YAMLOption{
		config.Permissive(),
		config.Expand(os.LookupEnv),
		config.Static(newDefaultConfig()),
	}

	if f, ferr := os.Stat(configFile); ferr == nil {
		log.Debug().
			Str("config file", configFile).
			Msg("using local config file")
		opts = append(opts, config.File(path.Join(f.Name())))
	}
	return config.NewYAML(opts...)
}
