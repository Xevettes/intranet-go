package main

import (
	"api/internal/config"
	"api/internal/gateways"
	"api/internal/server"
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	appCtx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("Failed load configurations", "error", err)
		os.Exit(1)
	}

	gatewayManager, err := gateways.NewManager(cfg)
	if err != nil {
		slog.Error("Failed to initializate gateway manager", "error", err)
		os.Exit(1)
	}

	srv := server.NewServer(gatewayManager, cfg)
	httpServer := &http.Server{
		Addr:    ":5555",
		Handler: srv.Router(),
	}

	go func() {
		slog.Info("Starting API Server", "port", ":5555")
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Servidor HTTP falhou", "error", err)
			cancel()
		}
	}()

	<-appCtx.Done()
	slog.Info("Shutdown...")

	// Lógica de Graceful Shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		slog.Error("Error in Shutdown", "error", err)
	}
	slog.Info("API Server Ended.")
}
