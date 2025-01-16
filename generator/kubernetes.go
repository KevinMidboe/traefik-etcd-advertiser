package generator

import (
	"fmt"
	"os"

	"github.com/kevinmidboe/traefik-etcd-advertiser/client/etcd"
	"gopkg.in/yaml.v3"
	"k8s.io/api/apps/v1"
)

func createPacket3(config *v1.Deployment) []etcd.EtcdPacket {
	blocks := []etcd.EtcdPacket{}

	fmt.Println("DockerToEtcd")
	fmt.Println("kube")
	fmt.Println(config.APIVersion)
	fmt.Println(*kubeConfig)

	return blocks
}

func KubernetesToEtcd(config *v1.Deployment, packetList *[]etcd.EtcdPacket) {
	items := createPacket3(config)

	*packetList = append(*packetList, items...)
}

func (filePath string) (*v1.Deployment, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cfg v1.Deployment
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
