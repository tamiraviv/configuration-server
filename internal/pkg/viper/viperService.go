package viper

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

func NewViperConfig() (*ViperConf, error) {
	viper.SetConfigName(defaultConfName) // name of config file (without extension)
	viper.AddConfigPath(defaultConfPath)
	viper.SetConfigType(defaultConfType)

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil { // Handle errors reading the config file
		return nil, errors.Wrap(err, "Failed to read configuration from file")
	}

	return &ViperConf{
		viper.GetViper(),
	}, nil
}

func (v *ViperConf) GetString(key string) string{
	return v.v.GetString(key)
}