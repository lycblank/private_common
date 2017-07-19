package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/coreos/etcd/client"
)

var etcdClient client.Client

func initEtcd(etcdaddrs []string) {
	cli, err := client.New(client.Config{
		Endpoints:               etcdaddrs,
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: 5 * time.Second,
	})
	if err != nil {
		fmt.Printf("new etcd client failed. error:%s\n", err)
		os.Exit(1)
	}
	etcdClient = cli
}

func GetEtcdClient() client.Client {
	return etcdClient
}
