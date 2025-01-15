package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/kevinmidboe/traefik-etcd-advertiser/client/etcd"
	"github.com/kevinmidboe/traefik-etcd-advertiser/config"
	"github.com/kevinmidboe/traefik-etcd-advertiser/converter"
	"github.com/kevinmidboe/traefik-etcd-advertiser/generator"
)

func getArgvFilename() string {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <path-to-yaml-file>\n", os.Args[0])
	}

	filename := os.Args[1]
	return filename
}

func main() {
	_, err := config.LoadConfig()
	if err != nil {
		log.Println("Error from config loader", err)
	}

	// setup etcd client
	// etcdManager, err := etcd.NewClient()
	if err != nil {
		panic(err)
	}

	var packets []etcd.EtcdPacket

	// parse traefik config from file
	filename := getArgvFilename()
	if strings.Contains(filename, "docker-compose.yml") {
		// build etcd packets from docker-compose config
		dockerConfig, err := generator.ParseDockerCompose(filename)
		if err != nil {
			log.Fatal(err)
		}

		// generator.PrintLabels(labels)
		fmt.Println("compose")
		log.Println(dockerConfig)
		// generator.TraefikToEtcd(traefikConfig, &packets)

	} else {
		// build etcd packets from traefik config
		traefikConfig, err := converter.TraefikFromYaml(filename)
		if err != nil {
			log.Fatalf("Error loading traefik YAML config file: %v\n", err)
		}

		generator.TraefikToEtcd(traefikConfig, &packets)
	}

	for _, packet := range packets {
		log.Println(packet)
		// etcdManager.Put(packet.Key, packet.Value)
	}
}
