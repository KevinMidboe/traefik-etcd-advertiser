package config

import (
	"fmt"
	"log/slog"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// Config contains environment variables.
type Config struct {
	EtcdEndpoint string `envconfig:"ETCD_ENDPOINTS"`
}

// LoadConfig reads environment variables, populates and returns Config.
func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		slog.Warn("No .env file found")
	}

	var c Config
	err := envconfig.Process("", &c)

	if len(c.EtcdEndpoint) < 1 {
		err = fmt.Errorf("missing variable ETCD_ENDPOINTS, not set")
	}

	return &c, err
}
