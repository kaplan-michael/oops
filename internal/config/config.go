package config

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// Config holds your application's configuration loaded from environment variables.
type Config struct {
	Template string
	Errors   string
	LogLevel string
	Port     string
}

// MustGetConfig loads configuration from environment variables and returns a Config pointer.
// It logs a fatal error if unmarshaling fails.
func MustGetConfig() *Config {
	viper.AutomaticEnv()

	// Set default values
	viper.SetDefault("TEMPLATE", "template.tmpl")
	viper.SetDefault("ERRORS", "errors.yaml")
	viper.SetDefault("LOGLEVEL", "info")
	viper.SetDefault("PORT", "8080")

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		zap.S().Fatalf("Unable to decode environment variables: %v", err)
	}

	zap.S().Info("Configuration loaded")
	return &cfg
}
