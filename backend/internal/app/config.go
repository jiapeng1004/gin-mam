package app

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Server        ServerConfig        `mapstructure:"server"`
	MySQL         MySQLConfig         `mapstructure:"mysql"`
	Redis         RedisConfig         `mapstructure:"redis"`
	Elasticsearch ElasticsearchConfig `mapstructure:"elasticsearch"`
	JWT           JWTConfig           `mapstructure:"jwt"`
	Log           LogConfig           `mapstructure:"log"`
}

type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type MySQLConfig struct {
	DSN string `mapstructure:"dsn"`
}

type RedisConfig struct {
	Addr string `mapstructure:"addr"`
	DB   int    `mapstructure:"db"`
}

type ElasticsearchConfig struct {
	Addresses []string `mapstructure:"addresses"`
}

type JWTConfig struct {
	Secret      string `mapstructure:"secret"`
	ExpireHours int    `mapstructure:"expire_hours"`
}

type LogConfig struct {
	Level string `mapstructure:"level"`
}

func LoadConfig(path string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(path)
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}
	if cfg.Server.Port == 0 {
		cfg.Server.Port = 8080
	}
	if cfg.Server.Mode == "" {
		cfg.Server.Mode = "debug"
	}
	if cfg.JWT.ExpireHours == 0 {
		cfg.JWT.ExpireHours = 24
	}
	return &cfg, nil
}
