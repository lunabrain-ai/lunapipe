package config

import (
	"github.com/lunabrain-ai/lunabrain/pkg/store/cache"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"os"
)

type Configurator interface {
	Configure(apiToken string) (string, error)
}

func (s *LocalConfigurator) Configure(apiToken string) (string, error) {
	c := NewDefaultConfig()
	c.OpenAI.APIKey = apiToken

	data, err := yaml.Marshal(c)
	if err != nil {
		return "", errors.Wrapf(err, "error marshalling YAML")
	}

	cfgFile, err := s.cache.GetFile(homeConfigFile)
	if err != nil {
		return "", errors.Wrapf(err, "error getting file")
	}

	err = os.WriteFile(cfgFile, data, 0644)
	if err != nil {
		return "", errors.Wrapf(err, "error writing YAML file")
	}
	return cfgFile, nil
}

type LocalConfigurator struct {
	cache cache.Cache
}

func NewConfigurator(cache cache.Cache) *LocalConfigurator {
	return &LocalConfigurator{
		cache: cache,
	}
}

var _ Configurator = (*LocalConfigurator)(nil)
