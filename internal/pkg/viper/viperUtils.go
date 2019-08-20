package viper

import "github.com/spf13/viper"

type ViperConf struct {
	v *viper.Viper
}

const (
	defaultConfPath = "configs/"
	defaultConfName = "config"
	defaultConfType = "yaml"
)

