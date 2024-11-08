package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	TCPServerPort    string   `mapstructure:"TCP_SERVER_PORT"`
	CertFile         string   `mapstructure:"CERT_FILE"`
	KeyFile          string   `mapstructure:"KEY_FILE"`
	AllowedAddresses []string `mapstructure:"ALLOWED_ADDRESSES"`
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
	if err == nil {
		// Split allowed addresses by comma if not empty
		allowedAddresses := viper.GetString("ALLOWED_ADDRESSES")
		if allowedAddresses != "" {
			config.AllowedAddresses = strings.Split(allowedAddresses, ",")
		}
	}

	return
}
