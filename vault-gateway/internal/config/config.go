package config

import (
	"fmt"
	"os"
	"time"
)

type Config struct {
	VaultSrvAddr      string
	VaultSrvRoleID    string
	VaultSrvSecretID  string
	VaultTimeout      time.Duration
	VaultGtwAuthToken string
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		VaultSrvAddr:      os.Getenv("VAULT_SERVER_ADDR"),
		VaultSrvRoleID:    os.Getenv("VAULT_SERVER_ROLE_ID"),
		VaultSrvSecretID:  os.Getenv("VAULT_SERVER_SECRET_ID"),
		VaultGtwAuthToken: os.Getenv("INTERNAL_API_AUTH_TOKEN"),
	}

	if cfg.VaultSrvAddr == "" {
		return nil, fmt.Errorf("VAULT_SERVER_ADDR environment variable is not set")
	}
	if cfg.VaultSrvRoleID == "" {
		return nil, fmt.Errorf("VAULT_SERVER_ROLE_ID environment variable is not set")
	}
	if cfg.VaultSrvSecretID == "" {
		return nil, fmt.Errorf("VAULT_SERVER_SECRET_ID environment variable is not set")
	}
	if cfg.VaultGtwAuthToken == "" {
		return nil, fmt.Errorf("INTERNAL_API_AUTH_TOKEN environment variable is not set")
	}

	cfg.VaultTimeout = 30 * time.Second // Default timeout

	return cfg, nil
}
