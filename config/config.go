package config

import (
	"github.com/spf13/viper"
	"fmt"
)

type Config struct {
	Telegram_Token string
	Db_Conn_Str string
}

var config Config

func C() *Config {
	return &config
}

func Init(file string) {
	viper.SetConfigName(file)
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Error in config file: %s", err))
	}

	viper.Unmarshal(&config)
}
