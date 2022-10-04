package config

import (
	"path/filepath"
	"strings"

	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
)

type Config struct {
	Environment string
	Debug       bool
	Host        string
	Port        string
	Db          map[string]string
	Broker      map[string]string
	Monitoring  map[string]string
}

func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

func (c *Config) IsDevelopment() bool {
	return c.Environment == "dev"
}

func GetConfig(configPath string) (*Config, error) {
	var config Config
	splits := strings.Split(filepath.Base(configPath), ".")
	viper.SetConfigName(filepath.Base(splits[0]))
	viper.AddConfigPath(filepath.Dir(configPath))

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
