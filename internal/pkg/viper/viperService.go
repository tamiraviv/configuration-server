package viper

import (
	"bytes"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

func NewViperConfig() (*ViperConf, error) {
	v := viper.GetViper()
	v.SetConfigName(defaultConfName) // name of config file (without extension)
	v.AddConfigPath(defaultConfPath)
	v.SetConfigType(defaultConfType)

	err := v.ReadInConfig() // Find and read the config file
	if err != nil {         // Handle errors reading the config file
		return nil, errors.Wrap(err, "Failed to read configuration from file")
	}

	return &ViperConf{
		v: v,
	}, nil
}

func (v *ViperConf) Get(key string) interface{} {
	return v.v.Get(key)
}

func (v *ViperConf) GetString(key string) string {
	return v.v.GetString(key)
}

func (v *ViperConf) Set(key string, value string) {
	v.v.Set(key, value)
}

func (v *ViperConf) FlushConfig() error {
	if err := viper.ReadConfig(bytes.NewReader(nil)); err != nil {
		return errors.Wrapf(err, "Failed to reset config")
	}
	return nil
}

func (v *ViperConf) SetConfig(m map[string]interface{}) error {
	b, err := json.Marshal(m)
	if err != nil {
		return errors.Wrapf(err, "Failed to marshal config map")
	}

	return v.v.ReadConfig(bytes.NewReader(b))
}
