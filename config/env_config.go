package config

import "github.com/kelseyhightower/envconfig"

type EnvConfig struct {
	SlackBotUserOAuthToken       string `envconfig:"SLACK_BOT_USER_OAUTH_TOKEN" required:"true"`
	GoogleApplicationCredentials string `envconfig:"GOOGLE_APPLICATION_CREDENTIALS" required:"true"`
	YamlConfigPath               string `envconfig:"YAML_CONFIG_PATH" default:"config.yaml"`
}

func LoadEnvConfig() (*EnvConfig, error) {
	c := &EnvConfig{}
	if err := envconfig.Process("", c); err != nil {
		return nil, err
	}
	return c, nil
}
