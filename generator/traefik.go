package generator

import (
	"fmt"
	"log/slog"
	"math"
	"strconv"

	"github.com/kevinmidboe/traefik-etcd-advertiser/client/etcd"
	"github.com/kevinmidboe/traefik-etcd-advertiser/converter"
	"github.com/traefik/traefik/v3/pkg/config/dynamic"
)

const traefikPrefix = "traefik"

func isKnownGenericType(value interface{}) bool {
	switch v := value.(type) {
	case string:
		return true
	case float64:
		return true
	case bool:
		return true
	default:
		slog.Debug(fmt.Sprintf("found unknown generic %s\n", v))
	}

	return false
}

func convertToGeneric(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case float64:
		return fmt.Sprintf("%d", int(math.Floor(v)))
	case bool:
		return strconv.FormatBool(v)
	}

	return "unknown type"
}

// recursively walks the JSON object and creates internal
// `etcdPackets` per leaf node, returns list of packets.
func createPacket(item interface{}, parentKey string) []etcd.EtcdPacket {
	blocks := []etcd.EtcdPacket{}

	switch itemD := item.(type) {
	// input is JSON object
	case map[string]interface{}:
		for key, value := range itemD {
			// check for generic value type vs nested object,
			// either create block or recursively call obj again

			if isKnownGenericType(value) {
				blocks = append(blocks, etcd.EtcdPacket{
					Key:   parentKey + "/" + key,
					Value: convertToGeneric(value),
				})
			} else {
				blocks = append(blocks, createPacket(itemD[key], fmt.Sprintf("%s/%s", parentKey, key))...)
			}
		}
	// input is JSON list
	case []interface{}:
		for i, item := range itemD {
			blocks = append(blocks, createPacket(item, parentKey+"/"+strconv.Itoa(i))...)
		}
	}

	return blocks
}

func TraefikToEtcd(config *dynamic.Configuration, packetList *[]etcd.EtcdPacket) {
	// always convert to json before converting to etcd
	data := converter.TraefikToJSON(config)

	// generate list of etcd commands
	items := createPacket(data, traefikPrefix)
	*packetList = append(*packetList, items...)
}

