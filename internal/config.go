package internal

import (
	"go.uber.org/config"
	"os"
	"path"
)

const configFile = ".lunapipe.yaml"

type Config struct {
	APIKey string `yaml:"api_key"`
}

type LogConfig struct {
	Level string `yaml:"level"`
}

type BaseConfig struct {
	OpenAI Config    `yaml:"openai"`
	Log    LogConfig `yaml:"log"`
}

func NewOpenAIConfig(provider config.Provider) (Config, error) {
	var c Config
	err := provider.Get("openai").Populate(&c)
	if err != nil {
		return Config{}, err
	}
	return c, nil
}

func NewLogConfig(provider config.Provider) (LogConfig, error) {
	var c LogConfig
	err := provider.Get("log").Populate(&c)
	if err != nil {
		return LogConfig{}, err
	}
	return c, nil
}

func newDefaultConfig() BaseConfig {
	return BaseConfig{
		OpenAI: Config{
			APIKey: "${OPENAI_API_KEY}",
		},
		Log: LogConfig{
			Level: "${LOG_LEVEL:info}",
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
		//log.Debug().
		//	Str("config file", configFile).
		//	Msg("using local config file")
		opts = append(opts, config.File(path.Join(f.Name())))
	}
	return config.NewYAML(opts...)
}
