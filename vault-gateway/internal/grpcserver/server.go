package grpcserver

import (
	"context"
	"fmt"
	"vault-gateway/internal/vault_client"
	"vault-gateway/proto/vault"

	"google.golang.org/protobuf/types/known/structpb"
)

type Server struct {
	vault.UnimplementedSecretServiceServer
	vaultClient *vault_client.VaultClient
}

func NewServer(vaultClient *vault_client.VaultClient) *Server {
	return &Server{
		vaultClient: vaultClient,
	}
}

func (s *Server) ReadSecret(ctx context.Context, req *vault.ReadSecretRequest) (*vault.ReadSecretResponse, error) {
	vaultSecret, err := s.vaultClient.ReadSecret(ctx, req.GetPath())
	if err != nil {
		return nil, err
	}

	dataStruct, err := structpb.NewStruct(vaultSecret.Data)
	if err != nil {
		return nil, err
	}

	return &vault.ReadSecretResponse{
		RequestId:     vaultSecret.RequestID,
		LeaseId:       vaultSecret.LeaseID,
		Renewable:     vaultSecret.Renewable,
		LeaseDuration: int32(vaultSecret.LeaseDuration),
		Data:          dataStruct,
	}, nil
}

func (s *Server) WriteSecret(ctx context.Context, req *vault.WriteSecretRequest) (*vault.WriteSecretResponse, error) {
	dataMap := req.GetData().AsMap()

	_, err := s.vaultClient.WriteSecret(ctx, req.GetPath(), dataMap)
	if err != nil {
		return nil, err
	}
	return &vault.WriteSecretResponse{
		Success: true,
	}, nil
}

func (s *Server) ListSecrets(ctx context.Context, req *vault.ListSecretsRequest) (*vault.ListSecretsResponse, error) {
	secret, err := s.vaultClient.List(ctx, req.GetPath())
	if err != nil {
		return nil, err
	}

	keysData, ok := secret.Data["keys"]
	if !ok {
		return &vault.ListSecretsResponse{Keys: []string{}}, nil
	}

	keysList, ok := keysData.([]interface{})
	if !ok {
		return nil, fmt.Errorf("formato inesperado para a lista de chaves do Vault")
	}

	var keys []string
	for _, key := range keysList {
		if strKey, ok := key.(string); ok {
			keys = append(keys, strKey)
		}
	}

	return &vault.ListSecretsResponse{Keys: keys}, nil
}
