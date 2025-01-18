package main

import (
	"log"
	"strings"

	"github.com/kevinmidboe/traefik-etcd-advertiser/client/etcd"
	"github.com/kevinmidboe/traefik-etcd-advertiser/config"
	"github.com/kevinmidboe/traefik-etcd-advertiser/converter"
	"github.com/kevinmidboe/traefik-etcd-advertiser/generator"
)

func main() {
	_, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error from config loader: %s", err)
	}

	filename, publish := config.ParseCli()

	// setup etcd client
	etcdManager, err := etcd.NewClient()
	if err != nil {
		panic(err)
	}

	var packets []etcd.EtcdPacket

	if strings.Contains(filename, "docker-compose.yml") {
		// build etcd packets from docker-compose config
		dockerConfig, err := generator.ParseDockerCompose(filename)
		if err != nil {
			log.Fatalf("Error loading docker YAML config file :%v\n", err)
		}

		generator.DockerToEtcd(dockerConfig, &packets)

	} else if strings.Contains(filename, "kubernetes") {
		// build etcd packets from kubernetes service resource
		kubeConfig, err := converter.ServiceToKubernetes(filename)
		if err != nil {
			log.Fatalf("Error loading traefik YAML config file: %v\n", err)
		}

		generator.KubernetesToEtcd(kubeConfig, &packets)
	} else {
		// build etcd packets from traefik config
		traefikConfig, err := converter.TraefikFromYaml(filename)
		if err != nil {
			log.Fatalf("Error loading traefik YAML config file: %v\n", err)
		}

		generator.TraefikToEtcd(traefikConfig, &packets)
	}

	etcd.RemoveDuplicatePackets(&packets)
	for _, packet := range packets {
		log.Println(packet)

		if *publish {
			etcdManager.Put(packet.Key, packet.Value)
		}
	}
}
