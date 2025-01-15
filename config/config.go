package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// Config contains environment variables.
type Config struct {
	EtcdEndpoint  string `envconfig:"ETCD_ENDPOINTS"`
}

// LoadConfig reads environment variables, populates and returns Config.
func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	var c Config

	err := envconfig.Process("", &c)

	return &c, err
}
