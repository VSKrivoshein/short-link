package store

import "github.com/spf13/viper"

type Config struct {
	DatabaseUrl string
}

func NewConfig() *Config {
	return &Config{DatabaseUrl: viper.GetString("database_url")}
}
