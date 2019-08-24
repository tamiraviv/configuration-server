package etcd

import (
	"configuration-server/internal/pkg/viper"

	"go.etcd.io/etcd/client"
)

type Etcd struct {
	client client.KeysAPI
	v      *viper.ViperConf
	dir    string
}

type WatcherResponse struct {
	Key   string
	Value string
	Err   error
}

const (
	etcdKey  = "etcd"
	etcdHost = etcdKey + ".host"
	etcdDir  = etcdKey + ".dir"
)
