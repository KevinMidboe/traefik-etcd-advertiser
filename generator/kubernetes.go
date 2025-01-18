package generator

import (
	"github.com/kevinmidboe/traefik-etcd-advertiser/client/etcd"
	"k8s.io/api/core/v1"
)

func KubernetesToEtcd(config *v1.Service, packetList *[]etcd.EtcdPacket) {
	for key, val := range config.ObjectMeta.Annotations {
	  *packetList = append(*packetList, etcd.EtcdPacket{
			Key: key,
			Value: val,
		})
	}
}
