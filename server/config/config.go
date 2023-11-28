package config

import (
	"context"
	"encoding/json"
	"time"

	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Config struct {
	LoggingLevel zerolog.Level      `mapstructure:"logging_level"`
	ListenAddr   string             `mapstructure:"listen_addr"`
	BackupsDir   string             `mapstructure:"backups_dir"`
	Backups      map[string]*Backup `mapstructure:"backups"`
}

type Backup struct {
	Interval   time.Duration `mapstructure:"interval"`
	Commands   []string      `mapstructure:"commands"`
	Subdir     string        `mapstructure:"subdir"`
	FilesLimit int           `mapstructure:"files_limit"`
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

func Log(ctx context.Context, config *Config) error {
	log := zlog.Ctx(ctx)

	msg, err := json.MarshalIndent(config, "config ", "  ")
	if err != nil {
		return err
	}

	log.Debug().Msg(string(msg))

	return nil
}
