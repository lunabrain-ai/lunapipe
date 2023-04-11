package internal

import (
	"go.uber.org/config"
	"os"
)

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
	return config.NewYAML(opts...)
}
