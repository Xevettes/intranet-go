package main

import (
	"context"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"zabbix-gateway/internal/auth"
	"zabbix-gateway/internal/config"
	"zabbix-gateway/internal/grpcserver"
	"zabbix-gateway/internal/zabbix_client"
	monitoring "zabbix-gateway/proto/zabbix"

	"google.golang.org/grpc"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	appCtx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("Erro ao carregar configurações", "error", err)
		os.Exit(1)
	}
	slog.Info("Configurações carregadas com sucesso.")

	zabbixClient, err := zabbix_client.NewClient(cfg.ZabbixAPIURL, cfg.ZabbixAPIToken)
	if err != nil {
		slog.Error("Falha ao inicializar cliente Zabbix", "error", err)
		os.Exit(1)
	}
	slog.Info("Cliente Zabbix inicializado com sucesso.")

	lis, err := net.Listen("tcp", ":5555")
	if err != nil {
		slog.Error("Falha ao escutar na porta", "error", err)
		os.Exit(1)
	}

	authInterceptor := auth.NewAuthInterceptor(cfg.GatewayAuthToken)
	gServer := grpc.NewServer(
		grpc.UnaryInterceptor(authInterceptor.Unary()),
	)

	zabbixGrpcServer := grpcserver.NewServer(zabbixClient)
	monitoring.RegisterMonitoringServiceServer(gServer, zabbixGrpcServer)

	go func() {
		slog.Info("Iniciando servidor gRPC do Zabbix Gateway", "port", ":5555")
		if err := gServer.Serve(lis); err != nil {
			slog.Error("Servidor gRPC foi encerrado inesperadamente", "error", err)
			cancel()
		}
	}()

	<-appCtx.Done()
	slog.Info("Sinal de shutdown recebido, iniciando desligamento gracioso...")

	gServer.GracefulStop()

	slog.Info("Aplicação Zabbix Gateway finalizada.")
}
