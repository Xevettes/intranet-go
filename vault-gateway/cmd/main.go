package main

import (
	"context"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"vault-gateway/internal/auth"
	"vault-gateway/internal/config"
	"vault-gateway/internal/grpcserver"
	"vault-gateway/internal/vault_client"
	"vault-gateway/proto/vault"

	"google.golang.org/grpc"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	appCtx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}
	slog.Info("Configuration loaded successfully", "vault_address", cfg.VaultSrvAddr)
	slog.Info("Vault Gateway starting...")

	authInterceptor := auth.NewAuthInterceptor(cfg.VaultGtwAuthToken)
	gServer := grpc.NewServer(grpc.UnaryInterceptor(authInterceptor.Unary()))
	slog.Info("gRPC server created")

	slog.Info("Creating Vault client...")
	vaultClient, err := vault_client.NewVaultClient(appCtx, cfg)
	if err != nil {
		slog.Error("Failed to create Vault client", "error", err)
		os.Exit(1)
	}
	slog.Info("Vault client created successfully")

	lis, err := net.Listen("tcp", ":5555")
	if err != nil {
		slog.Error("Failed to listen on port 5555", "error", err)
		os.Exit(1)
	}
	slog.Info("Listening on port 5555")
	server := grpcserver.NewServer(vaultClient)
	vault.RegisterSecretServiceServer(gServer, server)
	slog.Info("Vault service registered")

	go func() {
		slog.Info("Starting gRPC server...", "port", ":5555")
		if err := gServer.Serve(lis); err != nil {
			slog.Error("Failed to start gRPC server", "error", err)
			cancel()
		}
	}()
	slog.Info("gRPC server started successfully", "port", ":5555")

	<-appCtx.Done()
	slog.Info("Shutting down Vault Gateway...")
	gServer.GracefulStop()
	slog.Info("Vault Gateway stopped gracefully")
}
