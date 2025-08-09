package config

import (
	"fmt"
	"os"
	"strconv"
)

type VaultGatewayConfig struct {
	Enabled   bool
	Address   string
	RoleID    string
	SecretID  string
	AuthToken string
}

type ZabbixGatewayConfig struct {
	Enabled   bool
	APIURL    string
	APIToken  string
	AuthToken string
}

type APICentralConfig struct {
	RESTAuthToken string
}

type Config struct {
	VaultGateway             VaultGatewayConfig
	ZabbixGateway            ZabbixGatewayConfig
	APICentral               APICentralConfig
	GatewayInternalAuthToken string
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}

	cfg.APICentral.RESTAuthToken = os.Getenv("API_REST_AUTH_TOKEN")
	cfg.GatewayInternalAuthToken = os.Getenv("INTERNAL_API_AUTH_TOKEN")
	if cfg.GatewayInternalAuthToken == "" {
		return nil, fmt.Errorf("INTERNAL_API_AUTH_TOKEN is required")
	} else if cfg.APICentral.RESTAuthToken == "" {
		return nil, fmt.Errorf("API_REST_AUTH_TOKEN is required")
	}

	vaultEnabled, _ := strconv.ParseBool(os.Getenv("VAULT_GATEWAY_ENABLED"))
	cfg.VaultGateway.Enabled = vaultEnabled
	if cfg.VaultGateway.Enabled {
		cfg.VaultGateway.Address = os.Getenv("VAULT_GATEWAY_ADDRESS")

		if cfg.VaultGateway.Address == "" {
			return nil, fmt.Errorf("VAULT_GATEWAY_ADDRESS is required when VAULT_GATEWAY_ENABLED is true")
		}
	}

	zabbixEnabled, _ := strconv.ParseBool(os.Getenv("ZABBIX_GATEWAY_ENABLED"))
	cfg.ZabbixGateway.Enabled = zabbixEnabled
	if cfg.ZabbixGateway.Enabled {
		cfg.ZabbixGateway.APIURL = os.Getenv("ZABBIX_GATEWAY_API_URL")

		if cfg.ZabbixGateway.APIURL == "" {
			return nil, fmt.Errorf("ZABBIX_GATEWAY_API_URL is required when ZABBIX_ENABLED is true")
		}
	}

	return cfg, nil
}
