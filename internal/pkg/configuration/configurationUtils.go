package configuration

import (
	"github.com/spf13/viper"
)

type Configuration struct {
	v     *viper.Viper
	etcd  *Etcd
	hooks map[string]func(value string) error
}

const (
	defaultConfPath        = "configs/"
	defaultConfName        = "config"
	defaultConfType        = "yaml"
	configurationServerKey = "configurationServer"
	etcdKey                = "etcd"
	etcdHostKey            = etcdKey + ".host"
	etcdDirKey             = etcdKey + ".dir"
)
