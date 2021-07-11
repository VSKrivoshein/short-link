package apiserver

import (
	"github.com/VSKrivoshein/short-link/internal/app/store"
	"github.com/spf13/viper"
)

type Config struct {
	BindAddr string
	LogLevel string
	Store *store.Config
}

func InitConfig(path string) error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("apiserver")
	return viper.ReadInConfig()
}

func NewConfig() *Config {
	return &Config{
		BindAddr: viper.GetString("bind_adder"),
		LogLevel: viper.GetString("log_level"),
		Store: store.NewConfig(),
	}
}
