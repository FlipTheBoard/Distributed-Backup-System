package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	ServerAddr string `mapstructure:"server_addr"`
	BackupsDir string `mapstructure:"backups_dir"`
}

func ParseConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	var config Config
	if err = viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, err
}
