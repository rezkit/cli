package config

import "github.com/spf13/viper"

var config *viper.Viper

func init() {
	config = viper.New()

	config.SetConfigName("rezkit")
	config.SetConfigType("yaml")
	config.AddConfigPath("$HOME/.config")
}

func GetConfig() *viper.Viper {
	return config
}
