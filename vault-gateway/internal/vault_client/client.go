package vault_client

import (
	"context"
	"fmt"
	"log/slog"
	"vault-gateway/internal/config"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/api/auth/approle"
)

type VaultClient struct {
	vaultClient *api.Client
}

func NewVaultClient(ctx context.Context, cfg *config.Config) (*VaultClient, error) {
	clientConfig := &api.Config{Address: cfg.VaultSrvAddr, Timeout: cfg.VaultTimeout}

	client, err := api.NewClient(clientConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create Vault client: %w", err)
	}

	appRoleAuth, err := approle.NewAppRoleAuth(cfg.VaultSrvRoleID, &approle.SecretID{FromString: cfg.VaultSrvSecretID})
	if err != nil {
		return nil, fmt.Errorf("failed to create AppRole auth: %w", err)
	}

	authInfo, err := client.Auth().Login(ctx, appRoleAuth)
	if err != nil {
		return nil, fmt.Errorf("failed to login to Vault: %w", err)
	}
	slog.Info("Vault login successful")

	go setupTokenRenewal(ctx, client, authInfo)
	return &VaultClient{vaultClient: client}, nil
}

func setupTokenRenewal(ctx context.Context, client *api.Client, token *api.Secret) {
	watcher, err := client.NewLifetimeWatcher(&api.LifetimeWatcherInput{
		Secret:    token,
		Increment: 3600,
	})
	if err != nil {
		slog.Error("Failed to create lifetime watcher", "error", err)
		return
	}

	slog.Info("Starting Vault token renewal watcher")
	go watcher.Start()
	slog.Info("Vault token renewal watcher started")
	go func() {
		defer watcher.Stop()
		for {
			select {
			case err := <-watcher.DoneCh():
				if err != nil {
					slog.Error("Vault token renewal watcher stopped with error", "error", err)
				} else {
					slog.Info("Vault token renewal watcher stopped successfully")
				}
				return
			case renewal := <-watcher.RenewCh():
				if renewal != nil {
					slog.Info("Vault token renewed", "lease_id", renewal.Secret.LeaseID, "renewable", renewal.Secret.Renewable, "lease_duration", renewal.Secret.LeaseDuration)
				} else {
					slog.Info("Vault token renewed with no additional information")
				}
			case <-ctx.Done():
				slog.Info("Context done, stopping Vault token renewal watcher")
				return
			}
		}
	}()
}

func (vc *VaultClient) ReadSecret(ctx context.Context, path string) (*api.Secret, error) {
	secret, err := vc.vaultClient.Logical().Read(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read from Vault at path %s: %w", path, err)
	}
	if secret == nil {
		return nil, fmt.Errorf("no data found at path %s", path)
	}
	return secret, nil
}

func (vc *VaultClient) WriteSecret(ctx context.Context, path string, data map[string]interface{}) (*api.Secret, error) {
	secret, err := vc.vaultClient.Logical().Write(path, data)
	if err != nil {
		return nil, fmt.Errorf("failed to write to Vault at path '%s': %w", path, err)
	}
	if secret == nil {
		return nil, fmt.Errorf("no response from Vault after writing to path %s", path)
	}
	return secret, nil
}

func (vc *VaultClient) List(ctx context.Context, path string) (*api.Secret, error) {
	secret, err := vc.vaultClient.Logical().List(path)
	if err != nil {
		return nil, fmt.Errorf("failed to list at Vault path %s: %w", path, err)
	}
	if secret == nil {
		return nil, fmt.Errorf("no data found at Vault path %s", path)
	}
	return secret, nil
}

func (vc *VaultClient) Delete(ctx context.Context, path string) error {
	secret, err := vc.vaultClient.Logical().Delete(path)
	if err != nil {
		return fmt.Errorf("failed to delete at Vault path %s: %w", path, err)
	}
	if secret == nil {
		return fmt.Errorf("no response from Vault after deleting at path %s", path)
	}
	return nil
}
