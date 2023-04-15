package config

import (
	"github.com/lunabrain-ai/lunabrain/pkg/store/cache"
	"go.uber.org/config"
	"os"
	"path"
	"time"
)

const (
	localConfigFile = ".lunapipe.yaml"
	homeConfigFile  = "config.yaml"
)

type OpenAIConfig struct {
	APIKey  string        `yaml:"api_key"`
	Timeout time.Duration `yaml:"timeout"`
}

type LogConfig struct {
	Level string `yaml:"level"`
}

type BaseConfig struct {
	Cache  cache.Config `yaml:"cache"`
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

func NewDefaultConfig() BaseConfig {
	return BaseConfig{
		Cache: cache.Config{
			Name: ".lunapipe",
		},
		OpenAI: OpenAIConfig{
			APIKey:  "${OPENAI_API_KEY:\"\"}",
			Timeout: time.Minute * 5,
		},
		Log: LogConfig{
			Level: "${LOG_LEVEL:info}",
		},
	}
}

func NewConfigProvider(cache cache.Cache) (config.Provider, error) {
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
