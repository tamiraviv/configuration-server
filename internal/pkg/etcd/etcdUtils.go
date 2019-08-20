package etcd

import "go.etcd.io/etcd/client"

type Etcd struct {
	client client.KeysAPI
}

const (
	etcdKey = "etcd"
	etcdHost = etcdKey + ".host"
)
