package util

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	ChainID        int64  `mapstructure:"chainId"`
	GasLimit       uint64 `mapstructure:"gasLimit"`
	AccountAddress string `mapstructure:"accountAddress"`
	PK             string `mapstructure:"pk"`
	Telegram       struct {
		Token      string `mapstructure:"token"`
		ChatID     int64  `mapstructure:"chatId"`
		UserName   string `mapstructure:"userName"`
		WebHookUrl string `mapstructure:"webHookUrl"`
	} `mapstructure:"telegram"`
}

func GetConfigInfo() Config {
	var config Config

	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("config")

	if err := v.ReadInConfig(); err != nil {
		log.Printf("Error reading config file, %s", err)
	}

	if err := v.Unmarshal(&config); err != nil {
		log.Fatalf("Unable to unmarshal config into struct: %v", err)
	}

	return config
}
