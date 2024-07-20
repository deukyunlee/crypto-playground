package util

import (
	"github.com/spf13/viper"
	"log"
)

func GetViper() *viper.Viper {
	v := viper.New()
	v.SetConfigName("axs_staking_info")
	v.SetConfigType("yaml")
	v.AddConfigPath("./config")

	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	return v
}
