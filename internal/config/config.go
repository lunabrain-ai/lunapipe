package config

import (
	"go.uber.org/config"
	"os"
	"path"
	"time"
)

const configFile = ".lunapipe.yaml"

type OpenAIConfig struct {
	APIKey  string        `yaml:"api_key"`
	Timeout time.Duration `yaml:"timeout"`
}

type LogConfig struct {
	Level string `yaml:"level"`
}

type BaseConfig struct {
	OpenAI OpenAIConfig `yaml:"openai"`
	Log    LogConfig    `yaml:"log"`
}

func NewOpenAIConfig(provider config.Provider) (OpenAIConfig, error) {
	var c OpenAIConfig
	err := provider.Get("openai").Populate(&c)
	if err != nil {
		return OpenAIConfig{}, err
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
		OpenAI: OpenAIConfig{
			APIKey:  "${OPENAI_API_KEY}",
			Timeout: time.Minute * 5,
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
