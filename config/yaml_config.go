package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type YamlConfig struct {
	Calendars map[string]Calendar `yaml:"calendars"`
}

type Calendar struct {
	ID           string   `yaml:"id"`
	NotifyBefore int      `yaml:"notify_before" default:"10"` // minutes
	Channels     []string `yaml:"channels"`
}

func LoadYamlConfig(path string) (*YamlConfig, error) {
	c := &YamlConfig{}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	if err := yaml.NewDecoder(f).Decode(c); err != nil {
		return nil, err
	}
	return c, nil
}
