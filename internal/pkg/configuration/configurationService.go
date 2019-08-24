package configuration

import (
	"fmt"
	"time"

	"configuration-server/internal/pkg/etcd"
	"configuration-server/internal/pkg/viper"

	"github.com/pkg/errors"
)

func NewConfigurationService() (*Configuration, error) {
	v, err := viper.NewViperConfig()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create Viper Config")
	}

	hooks := make(map[string]func(value string) error)

	configServer := v.GetString(configurationServerKey)
	if configServer == etcdKey {
		e, err := etcd.NewEtcd(v)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to create Etcd config server")
		}

		if err := v.FlushConfig(); err != nil {
			return nil, errors.Wrap(err, "Failed to flush config")
		}

		m, err := e.GetConfig()
		if err != nil {
			return nil, errors.Wrap(err, "Failed to get config from etcd")
		}

		if err := v.SetConfig(m); err != nil {
			return nil, errors.Wrap(err, "Failed to set config")
		}

		conf := Configuration{
			etcd:  e,
			viper: v,
			hooks: hooks,
		}
		go func() {
			for {
				watchRes, done := e.WatchConfig()
				select {
				case err := <-done:
					fmt.Println("Failed to watch remote config changes:", err)
					fmt.Println("Retrying watch remote config changes in 10 sec...")
					time.Sleep(10 * time.Second)
				case res := <-watchRes:
					conf.viper.Set(res.Key, res.Value)
					if err := conf.invokeHook(res.Key, res.Value); err != nil {
						fmt.Printf("Warning: Failed to invoke hook on key (%s) with value (%s): (%s)\n", res.Key, res.Value, err)
					}
				}
			}
		}()

		return &conf, nil
	}

	return &Configuration{
		viper: v,
		hooks: hooks,
	}, nil
}

func (c *Configuration) GetString(key string) string {
	return c.viper.GetString(key)
}

func (c *Configuration) Get(key string) interface{} {
	return c.viper.Get(key)
}

func (c *Configuration) Set(key string, value string) error {
	if c.etcd != nil {
		if err := c.etcd.Set(key, value); err != nil {
			return errors.Wrap(err, "Failed to set key (%s) value (%s) in configuration service")
		}
		return nil
	}

	c.viper.Set(key, value)
	if err := c.invokeHook(key, value); err != nil {
		fmt.Printf("Warning: Failed to invoke hook on key (%s) with value (%s): (%s)\n", key, value, err)
	}
	return nil
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
		return errors.Wrapf(err, "Failed to run hook on value (%s)", value)
	}

	return nil
}
