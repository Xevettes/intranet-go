package config

import (
	"fmt"
	"os"
)

type Config struct {
	NetboxAPIURL     string
	NetboxAPIToken   string
	GatewayAuthToken string
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		NetboxAPIURL:     os.Getenv("NETBOX_API_URL"),
		NetboxAPIToken:   os.Getenv("NETBOX_API_TOKEN"),
		GatewayAuthToken: os.Getenv("INTERNAL_API_AUTH_TOKEN"),
	}

	if cfg.NetboxAPIURL == "" {
		return nil, fmt.Errorf("variável de ambiente obrigatória não definida: NETBOX_API_URL")
	}
	if cfg.NetboxAPIToken == "" {
		return nil, fmt.Errorf("variável de ambiente obrigatória não definida: NETBOX_API_TOKEN")
	}
	if cfg.GatewayAuthToken == "" {
		return nil, fmt.Errorf("variável de ambiente obrigatória não definida: GATEWAY_AUTH_TOKEN")
	}

	return cfg, nil
}
