package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"github.com/davecgh/go-spew/spew"

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
		log.Fatalf("Error from config loader: %s", err)
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
			log.Fatalf("Error loading docker YAML config file :%v\n", err)
		}

		generator.DockerToEtcd(dockerConfig, &packets)

	} else if strings.Contains(filename, "deployment") {
		kubeConfig, err := generator.KubernetesToEtcd(filename)
		if err != nil {
			log.Fatalf("Error loading traefik YAML config file: %v\n", err)
		}

		fmt.Println("kube")
		fmt.Println(*kubeConfig)
		fmt.Println(*kubeConfig.Spec.Replicas)
		fmt.Printf("as: %+v\n", kubeConfig.Spec.Selector.MatchLabels["app"])
		spew.Dump(*kubeConfig.Sepc.Selector)
		
		fmt.Println(kubeConfig.ObjectMeta.Name)
		fmt.Println(kubeConfig.GetObjectMeta())
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
		// etcdManager.Put(packet.Key, packet.Value)
	}
}
