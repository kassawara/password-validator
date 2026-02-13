package config

import (
	"github.com/spf13/viper"
)

var C = &AppConfig{}

type AppConfig struct {
	AppName        string `mapstructure:"app_name"`
	AppVersion     string `mapstructure:"app_version"`
	Environment    string `mapstructure:"environment"`
	LoggingLevel   string `mapstructure:"logging_level"`
	HttpServerPort string `mapstructure:"http_server_port"`
	ServerTimeout  string `mapstructure:"server_timeout"`
}

func Load() error {
	v := viper.New()
	v.BindEnv("app_name")
	v.BindEnv("app_version")
	v.BindEnv("environment")
	v.SetDefault("logging_level", "info")
	v.SetDefault("http_server_port", "8080")
	v.SetDefault("server_timeout", "10")
	v.BindEnv("logging_level")
	v.BindEnv("http_server_port")
	v.BindEnv("server_timeout")

	v.AutomaticEnv()
	if err := v.Unmarshal(C); err != nil {
		return err
	}

	return nil
}
