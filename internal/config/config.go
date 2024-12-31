package config

import (
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

// ResticControllerConfig the collection of all configs
type AppConfig struct {
	Log          LogConfig     `mapstructure:"log"`
	Repositories []*Repository `mapstructure:"repositories" validate:"required,dive"`
}

// ReloadConfig (re-)loads and validates the config from the config file
func ReloadConfig(configFile string) (*AppConfig, error) {
	config := AppConfig{}

	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	err := viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	for _, repository := range config.Repositories {
		if repository.PasswordFile != "" {
			content, err := os.ReadFile(repository.PasswordFile)
			if err != nil {
				return nil, err
			}

			repository.Password = string(content)
		}

		if repository.Env == nil {
			repository.Env = make(map[string]string)
		}

		for k, v := range repository.EnvFromFile {
			content, err := os.ReadFile(v)
			if err != nil {
				return nil, err
			}
			repository.Env[k] = string(content)
		}
	}

	if err := validateConfig(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

// validateConfig validates the config file
func validateConfig(config *AppConfig) error {
	validate := validator.New()

	return validate.Struct(config)
}
