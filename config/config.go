package config

import "github.com/AlekseiKromski/at-socket-server/core"

type DbConfig struct {
	Host     string
	Username string
	Password string
	Database string
}

func NewDbConfig(host, username, password, database string) *DbConfig {
	return &DbConfig{
		Host:     host,
		Username: username,
		Password: password,
		Database: database,
	}
}

type Config struct {
	AppConfig *core.Config
	DbConfig  *DbConfig
}

func NewConfig(ap *core.Config, dc *DbConfig) *Config {
	return &Config{
		AppConfig: ap,
		DbConfig:  dc,
	}
}
