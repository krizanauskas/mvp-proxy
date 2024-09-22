package appconfig

import (
	"fmt"

	"github.com/spf13/viper"
)

type ProxyServerConfig struct {
	Port                  string `mapstructure:"port"`
	MaxRequestDurationSec int    `mapstructure:"max_request_duration_sec"`
	AllowedDataMB         int    `mapstructure:"allowed_data_mb"`
}

type StatusServerConfig struct {
	Port string `mapstructure:"port"`
}

type AppConfig struct {
	ProxyServer  ProxyServerConfig  `mapstructure:"proxy_server"`
	StatusServer StatusServerConfig `mapstructure:"status_server"`
}

func Init(env string) (*AppConfig, error) {
	var conf AppConfig

	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")
	viper.SetConfigName(fmt.Sprintf("config-%s", env))

	if err := viper.ReadInConfig(); err != nil {
		return &conf, fmt.Errorf("failed to read config: %w", err)
	}

	if err := viper.Unmarshal(&conf); err != nil {
		return &conf, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &conf, nil
}
