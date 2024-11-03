package util

import (
	"github.com/spf13/viper"
)

type Config struct {
	ChainID  int64  `mapstructure:"chainId"`
	GasLimit uint64 `mapstructure:"gasLimit"`
	PK       string `mapstructure:"pk"`
	Telegram struct {
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
		logger.Errorf("Error reading config file, %s", err)
	}

	if err := v.Unmarshal(&config); err != nil {
		logger.Errorf("Unable to unmarshal config into struct: %v", err)
	}

	return config
}
