package apiserver

import "github.com/spf13/viper"

type Config struct {
	BindAddr string
}

func InitConfig(path string) error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("apiserver")
	return viper.ReadInConfig()
}

func NewConfig() *Config {
	return &Config{BindAddr: viper.GetString("bind_adder")}
}
