package main

import (
	"context"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"netbox-gateway/internal/auth"
	"netbox-gateway/internal/config"
	"netbox-gateway/internal/grpcserver"
	"netbox-gateway/internal/netbox_client"
	"netbox-gateway/proto/dcim_proto"
	"netbox-gateway/proto/ipam"
	"netbox-gateway/proto/organization"
	"netbox-gateway/proto/virtualization"

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
	slog.Info("Configuration loaded successfully", "netbox_url", cfg.NetboxAPIURL)

	netboxClient, err := netbox_client.NewClient(cfg.NetboxAPIURL, cfg.NetboxAPIToken)
	if err != nil {
		slog.Error("Failed to create NetBox client", "error", err)
		os.Exit(1)
	}
	slog.Info("NetBox client created successfully")

	lis, err := net.Listen("tcp", ":5555")
	if err != nil {
		slog.Error("Failed to listen on port 5555", "error", err)
		os.Exit(1)
	}
	authInterceptor := auth.NewAuthInterceptor(cfg.GatewayAuthToken)
	gServer := grpc.NewServer(grpc.UnaryInterceptor(authInterceptor.Unary()))

	server := grpcserver.NewServer(netboxClient)

	// Registrar todos os serviços
	organization.RegisterOrganizationServiceServer(gServer, server)
	dcim_proto.RegisterDcimServiceServer(gServer, server)
	ipam.RegisterIpamServiceServer(gServer, server)
	virtualization.RegisterVirtualizationServiceServer(gServer, server)

	go func() {
		slog.Info("Starting NetBox Gateway gRPC server", "port", ":5555")
		if err := gServer.Serve(lis); err != nil {
			slog.Error("gRPC server failed", "error", err)
			cancel()
		}
	}()

	<-appCtx.Done()
	slog.Info("Shutting down NetBox Gateway...")
	gServer.GracefulStop()
	slog.Info("NetBox Gateway stopped gracefully")
}
