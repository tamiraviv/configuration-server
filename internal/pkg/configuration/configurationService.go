package configuration

import (
	"configuration-server/internal/pkg/etcd"
	"configuration-server/internal/pkg/viper"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func NewConfigurationService() (*Configuration, error) {
	v, err := viper.NewViperConfig()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create Viper Config")
	}

	e, err := etcd.NewEtcd(v)
	if err != nil {
		logrus.Warningln("Could not create etcd client:", err)
	}

	return &Configuration{
		etcd: e,
		viper: v,
	}, nil
}

func (c *Configuration) Get(key string) (string, error){
	if c.etcd != nil {
		return c.etcd.Get(key)
	}

	return c.viper.GetString(key), nil
}
