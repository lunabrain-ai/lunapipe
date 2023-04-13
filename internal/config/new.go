package config

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"os"
)

func writeToYamlFile(c BaseConfig) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return errors.Wrapf(err, "error marshalling YAML")
	}

	err = os.WriteFile(".lunapipe.yaml", data, 0644)
	if err != nil {
		return errors.Wrapf(err, "error writing YAML file")
	}
	println("Wrote .lunapipe.yaml file")
	return nil
}

func NewConfigurator(apiToken string) error {
	c := BaseConfig{
		OpenAI: Config{
			APIKey: apiToken,
		},
	}
	return writeToYamlFile(c)
}
