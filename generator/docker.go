package generator

import (
	"context"
	"fmt"
	"strings"

	"github.com/compose-spec/compose-go/v2/cli"
	"github.com/compose-spec/compose-go/v2/types"
	"github.com/kevinmidboe/traefik-etcd-advertiser/client/etcd"
)

func dockerLabelToEtcdKey(key string) string {
	return strings.ReplaceAll(key, ".", "/")
}

func ParseDockerCompose(composeFilePath string) (*types.Project, error) {
	ctx := context.Background()

	options, err := cli.NewProjectOptions(
		[]string{composeFilePath},
	)
	if err != nil {
		return nil, err
	}

	project, err := options.LoadProject(ctx)
	if err != nil {
		return nil, err
	}

	return project, nil
}

func createPacket2(config *types.Project) []etcd.EtcdPacket {
	blocks := []etcd.EtcdPacket{}

	fmt.Println("DockerToEtcd")
	if config.Services == nil || len(config.Services) < 1 {
		fmt.Println("services not found - skipping")
		return blocks
	}

	for serviceName, _ := range config.Services {
		if config.Services[serviceName].Labels == nil {
			fmt.Println("sevice, but no labels found - continuing")
			continue
		}

		for key, val := range config.Services[serviceName].Labels {
			blocks = append(blocks, etcd.EtcdPacket{
				Key:   dockerLabelToEtcdKey(key),
				Value: val,
			})
		}
	}

	return blocks
}

func DockerToEtcd(config *types.Project, packetList *[]etcd.EtcdPacket) {
	items := createPacket2(config)

	*packetList = append(*packetList, items...)
}
