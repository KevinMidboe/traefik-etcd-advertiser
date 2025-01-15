package converter

import (
	"os"

	"github.com/traefik/traefik/v3/pkg/config/dynamic"
	"gopkg.in/yaml.v3"
)

func TraefikFromYaml(filePath string) (*dynamic.Configuration, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cfg dynamic.Configuration
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func TraefikToYaml(config *dynamic.Configuration) string {
	return ""
}
