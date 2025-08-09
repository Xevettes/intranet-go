package zabbix_client

import (
	authinterceptor "api/internal/grpcclients/auth_interceptor"
	"api/proto/zabbix"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// NewZabbixClient cria uma conexão e retorna um cliente gRPC para o Zabbix Gateway.
func NewZabbixClient(gatewayAddress string, authToken string) (monitoring.MonitoringServiceClient, error) {
	authInterceptor := authinterceptor.NewAuthInterceptor(authToken)
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(authInterceptor.Unary()),
	}

	slog.Info("Connecting to Zabbix service", "address", gatewayAddress)
	conn, err := grpc.NewClient(gatewayAddress, opts...)
	if err != nil {
		slog.Error("Failed to connect to Zabbix service", "error", err)
		return nil, err
	}

	client := monitoring.NewMonitoringServiceClient(conn)
	slog.Info("Successfully connected to Zabbix service", "address", gatewayAddress)
	return client, nil
}
