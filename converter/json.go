package converter

import (
	"encoding/json"
	"log"

	"github.com/traefik/traefik/v3/pkg/config/dynamic"
)

func TraefikToJSON(config *dynamic.Configuration) map[string]interface{} {
	var data map[string]interface{}
	jsonData, _ := json.Marshal(config)
	if err := json.Unmarshal(jsonData, &data); err != nil {
		log.Printf("failed to unmarshal JSON: %w", err)
	}
	
	return data
}
