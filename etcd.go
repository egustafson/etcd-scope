package main

import (
	"time"

	etcd "go.etcd.io/etcd/client/v3" // <-- v3.5
)

type Etcd3Config struct {
	Endpoint string
}

func newEtcd3Client() (c *etcd.Client, err error) {

	return etcd.New(etcd.Config{
		Endpoints:        []string{Config.Endpoint},
		AutoSyncInterval: 0, // <-- disables auto-sync
		DialTimeout:      1 * time.Second,
	})
}
