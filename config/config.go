package config

import (
	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
)

type Config struct {
	Debug  bool
	Host   string
	Port   string
	Db     map[string]string
	Broker map[string]string
}

func GetConfig(configPath string) (*Config, error) {
	var config Config
	viper.SetConfigName(".config")
	viper.AddConfigPath(".")
	viper.AddConfigPath(configPath)
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Unable to get config %s", err)
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("Unable to decode config into struct: %v", err)
	}
	return &config, nil
}
