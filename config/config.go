package config

import "github.com/AlekseiKromski/at-socket-server/core"

type ApiConfig struct {
	JWTSecret []byte
}

func NewApiConfig(jwtSecret string) *ApiConfig {
	return &ApiConfig{
		JWTSecret: []byte(jwtSecret),
	}
}

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
	ApiConfig *ApiConfig
}

func NewConfig(ap *core.Config, dc *DbConfig, ac *ApiConfig) *Config {
	return &Config{
		AppConfig: ap,
		DbConfig:  dc,
		ApiConfig: ac,
	}
}
