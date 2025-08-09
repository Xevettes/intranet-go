package vault

import (
	authinterceptor "api/internal/grpcclients/auth_interceptor"
	"api/proto/vault"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewClient(gatewayAddress string, authToken string) (vault.SecretServiceClient, error) {
	authInterceptor := authinterceptor.NewAuthInterceptor(authToken)
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(authInterceptor.Unary()),
	}

	slog.Info("Connecting to Vault service", "address", gatewayAddress)
	conn, err := grpc.NewClient(gatewayAddress, opts...)
	if err != nil {
		slog.Error("Failed to connect to Vault service", "error", err)
		return nil, err
	}

	slog.Info("Successfully connected to Vault service", "address", gatewayAddress)
	client := vault.NewSecretServiceClient(conn)
	slog.Info("Connected to Vault service", "address", gatewayAddress)
	return client, nil
}
