package gateways

import (
	"api/internal/config"
	"api/internal/grpcclients/vault"
	zabbix_client "api/internal/grpcclients/zabbix"
	vault_proto "api/proto/vault"
	zabbix_proto "api/proto/zabbix"
	"log/slog"
)

type Manager struct {
	VaultClient  vault_proto.SecretServiceClient
	ZabbixClient zabbix_proto.MonitoringServiceClient
}

func NewManager(cfg *config.Config) (*Manager, error) {
	manager := &Manager{}
	var err error

	if cfg.VaultGateway.Enabled {
		manager.VaultClient, err = vault.NewClient(cfg.VaultGateway.Address, cfg.GatewayInternalAuthToken)
		if err != nil {
			slog.Error("Failed to create Vault client", "error", err)
			return nil, err
		}
		slog.Info("Vault client created successfully")
	}
	if cfg.ZabbixGateway.Enabled {
		manager.ZabbixClient, err = zabbix_client.NewZabbixClient(cfg.ZabbixGateway.APIURL, cfg.GatewayInternalAuthToken)
		if err != nil {
			slog.Error("Failed to create Zabbix client", "error", err)
			return nil, err
		}
		slog.Info("Zabbix client initialized successfully")
	}

	return manager, nil
}
