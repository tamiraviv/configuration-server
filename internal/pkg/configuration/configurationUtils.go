package configuration

import (
	"configuration-server/internal/pkg/etcd"
	"configuration-server/internal/pkg/viper"
)

type Configuration struct {
	viper *viper.ViperConf
	etcd  *etcd.Etcd
	hooks map[string]func(value string) error
}

const (
	configurationServerKey = "configurationServer"
	etcdKey                = "etcd"
)
