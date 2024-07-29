package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	ServerPort string `mapstructure:"SERVER_PORT"`
	CertFile   string `mapstructure:"CERT_FILE"`
	KeyFile    string `mapstructure:"KEY_FILE"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		log.Printf("Error reading config file, %s", err)
	}

	err = viper.Unmarshal(&config)
	return
}
