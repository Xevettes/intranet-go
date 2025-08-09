package config

import (
	"fmt"
	"os"
)

type Config struct {
	ZabbixAPIURL     string
	ZabbixAPIToken   string
	GatewayAuthToken string
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		ZabbixAPIURL:     os.Getenv("ZABBIX_API_URL"),
		ZabbixAPIToken:   os.Getenv("ZABBIX_API_TOKEN"),
		GatewayAuthToken: os.Getenv("INTERNAL_API_AUTH_TOKEN"),
	}
	if cfg.ZabbixAPIURL == "" {
		return nil, fmt.Errorf("variável de ambiente obrigatória não definida: ZABBIX_API_URL")
	}
	if cfg.ZabbixAPIToken == "" {
		return nil, fmt.Errorf("variável de ambiente obrigatória não definida: ZABBIX_API_TOKEN")
	}
	if cfg.GatewayAuthToken == "" {
		return nil, fmt.Errorf("variável de ambiente obrigatória não definida: INTERNAL_API_AUTH_TOKEN")
	}
	return cfg, nil
}
