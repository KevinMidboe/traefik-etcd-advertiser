package etcd

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
	etcd "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc/connectivity"
)

type EtcdManager struct {
	client *etcd.Client
}

type EtcdPacket struct {
	Key   string
	Value string
}

func NewClient() (*EtcdManager, error) {
	fmt.Println("setting up etcd client")

	endpoints := []string{"localhost:2379"}
	if value, exists := os.LookupEnv("ETCD_ENDPOINTS"); exists {
		endpoints = strings.Split(value, ",")
	}

	client, err := etcd.New(etcd.Config{
		Endpoints:   endpoints,
		DialTimeout: 30 * time.Second,
	})
	if err != nil {
		return nil, errors.Wrap(err, "could not create etcd manager")
	}

	mgr := &EtcdManager{client}

	err = mgr.runWithTimeout(func(e *EtcdManager, ctx context.Context) error {
		etcdStatus, err := e.client.Status(ctx, e.client.Endpoints()[0])
		if err != nil {
			return fmt.Errorf("could not get etcd status: %v", err)
		}
		if len(etcdStatus.Errors) > 0 {
			return fmt.Errorf("etcd server has errors: %v", etcdStatus.Errors)
		}

		connState := e.client.ActiveConnection().GetState()
		if connState != connectivity.Ready && connState != connectivity.Idle {
			return fmt.Errorf("etcd connection is in unexpecetd state: %s",
				e.client.ActiveConnection().GetState())
		}
		return err
	})

	if err != nil {
		return nil, err
	}

	return mgr, nil
}

func (e *EtcdManager) runWithTimeout(runnable func(e *EtcdManager, ctx context.Context) error) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return runnable(e, ctx)
}

func runWithTimeoutReturning[R any](
	e *EtcdManager,
	runnable func(e *EtcdManager, ctx context.Context) (R, error),
) (R, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return runnable(e, ctx)
}

func (e *EtcdManager) Put(key string, val string, opts ...etcd.OpOption) error {
	log.Println("etcd: Putting key", key, val)
	return e.runWithTimeout(func(e *EtcdManager, ctx context.Context) error {
		_, err := e.client.Put(ctx, key, val, opts...)
		return err
	})
}

func (e *EtcdManager) Get(key string, opts ...etcd.OpOption) (string, error) {
	log.Println("etcd: Getting key", key)
	return runWithTimeoutReturning(e, func(e *EtcdManager, ctx context.Context) (string, error) {
		resp, err := e.client.Get(ctx, key, opts...)
		if err != nil {
			return "", nil
		}
		return string(resp.Kvs[0].Value), nil
	})
}

func (e *EtcdManager) Remove(key string) error {
	log.Println("etcd: Removing key", key)
	return e.runWithTimeout(func(e *EtcdManager, ctx context.Context) error {
		_, err := e.client.Delete(ctx, key)
		return err
	})
}

func RemoveDuplicatePackets(packets *[]EtcdPacket) {
	keys := make(map[string]bool)
	list := []EtcdPacket{}
	keysRm := []string{}

	for _, entry := range *packets {
		if _, ok := keys[entry.Key]; !ok {
			keys[entry.Key] = true
			list = append(list, entry)
		} else {
			keysRm = append(keysRm, entry.Key)
		}
	}

	if len(keysRm) > 0 {
		slog.Warn(fmt.Sprintf("keys squashed: %s", keysRm))
	}

	*packets = list
}
