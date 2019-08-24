package configuration

import (
	"bytes"
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

func NewConfigurationService() (*Configuration, error) {
	v := viper.GetViper()
	v.SetConfigName(defaultConfName) // name of config file (without extension)
	v.AddConfigPath(defaultConfPath)
	v.SetConfigType(defaultConfType)

	err := v.ReadInConfig() // Find and read the config file
	if err != nil {         // Handle errors reading the config file
		return nil, errors.Wrap(err, "Failed to read configuration from file")
	}

	hooks := make(map[string]func(value string) error)
	configServer := v.GetString(configurationServerKey)
	if configServer == etcdKey {
		etcdDir := v.GetString(etcdDirKey)
		etcdHost := v.GetString(etcdHostKey)
		if err := v.AddRemoteProvider(configServer, etcdHost, etcdDir); err != nil {
			return nil, errors.Wrapf(err, "Failed to add remote provider (%s)", etcdKey)
		}

		etcdClient, err := NewEtcd(etcdHost, etcdDir)
		if err != nil {
			return nil, errors.Wrapf(err, "Failed to connect to remote provider (%s)", etcdKey)
		}

		if err := viper.ReadConfig(bytes.NewReader(nil)); err != nil {
			return nil, errors.Wrapf(err, "Failed to reset config")
		}

		viper.RemoteConfig = etcdClient
		if err := viper.ReadRemoteConfig(); err != nil {
			return nil, errors.Wrapf(err, "Failed to read from remote provider (%s)", etcdKey)
		}

		go func() {
			if err := v.WatchRemoteConfigOnChannel(); err != nil {
				fmt.Println("Failed to watch remote config changes:", err)
			}
			/*for {
				if err := v.WatchRemoteConfigOnChannel(); err != nil {
					fmt.Println("Failed to watch remote config changes:", err)
				}
				fmt.Println("Retrying watch remote config changes in 10 sec...")
				time.Sleep(10 * time.Second)
			}*/
		}()

		return &Configuration{
			v:     viper.GetViper(),
			etcd:  etcdClient,
			hooks: hooks,
		}, nil

	}

	return &Configuration{
		v:     viper.GetViper(),
		hooks: hooks,
	}, nil
}

func (c *Configuration) GetString(key string) string {
	return c.v.GetString(key)
}

func (c *Configuration) Set(key string, value string) error {
	if c.etcd != nil {
		if err := c.etcd.Set(key, value); err != nil {
			return errors.Wrap(err, "Failed to set key (%s) value (%s) in configuration service")
		}
	} else {
		c.v.Set(key, value)
	}

	return c.invokeHook(key, value)
}

func (c *Configuration) RegisterHook(key string, f func(value string) error) {
	c.hooks[key] = f
}

func (c *Configuration) invokeHook(key string, value string) error {
	f, ok := c.hooks[key]
	if !ok {
		return nil
	}

	if err := f(value); err != nil {
		return errors.Wrapf(err, "Failed to invoke hooks on key (%s) with value (%s)", key, value)
	}

	return nil
}
