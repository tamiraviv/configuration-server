package configuration

import (
	"configuration-server/internal/pkg/etcd"
	"configuration-server/internal/pkg/viper"
)

type Configuration struct {
	viper *viper.ViperConf
	etcd *etcd.Etcd
}
