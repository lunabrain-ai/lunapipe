package openai

import (
	"github.com/sashabaranov/go-openai"
	"go.uber.org/config"
	"time"
)

const defaultModel = openai.GPT3Dot5Turbo

type Config struct {
	APIKey  string        `yaml:"api_key"`
	Timeout time.Duration `yaml:"timeout"`
	Model   string        `yaml:"model"`
}

func NewDefaultConfig() Config {
	return Config{
		APIKey:  "${OPENAI_API_KEY:\"\"}",
		Timeout: time.Minute * 5,
		Model:   defaultModel,
	}
}

func NewConfig(provider config.Provider) (Config, error) {
	var c Config
	err := provider.Get("openai").Populate(&c)
	if err != nil {
		return Config{}, err
	}
	return c, nil
}
