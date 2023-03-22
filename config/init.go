package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("ERROR_READING_CONFIG_FILE", err)
		//logs.Error("ERROR_READING_CONFIG_FILE")
		return
	}
	fmt.Println("SUCCESS_READING_CONFIG_FILE")
}

func GetEnv(key, defaultValue string) string {
	getString := viper.GetString(key)
	if getString == "" {
		return defaultValue
	}
	return getString
}
