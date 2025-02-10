package main

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"url-shortener/cmd/url-shortener/server/grpcserver"
	"url-shortener/cmd/url-shortener/server/httpserver"
	"url-shortener/internal/config"
	"url-shortener/internal/logger"
	"url-shortener/internal/service"
	"url-shortener/internal/storage"
)

const (
	ServersNumber    = 2
	shutdownDuration = 10 * time.Second
)

func main() {
	cfg := config.MustLoadConfig()
	log := logger.NewLogger(cfg.Log.Level)
	db, err := storage.NewStorage(&cfg.Storage, log)
	if err != nil {
		log.Error("Failed to initialize storage: " + err.Error())
		os.Exit(1)
	}
	log.Info("Initialized storage")

	shortener := service.NewShortener(db, log)
	httpServer, grpcServer, lis := initializeServers(cfg, shortener, log)
	defer func(lis net.Listener) {
		_ = lis.Close()
	}(lis)

	runServers(httpServer, grpcServer, lis, log)
}

func initializeServers(cfg *config.Config, shortener *service.Shortener, log *zap.Logger) (*http.Server, *grpc.Server, net.Listener) {
	httpServer := httpserver.NewHTTPServer(cfg.Server, shortener, log)
	log.Info(fmt.Sprintf("Starting HTTP server on %s", httpServer.Addr))

	lis, err := net.Listen("tcp", cfg.Server.GRPCPort)
	if err != nil {
		log.Error("Failed to start gRPC listener: " + err.Error())
		os.Exit(1)
	}

	grpcServer := grpcserver.New(cfg.Server, shortener, log)
	log.Info(fmt.Sprintf("Starting gRPC server on port %s", cfg.Server.GRPCPort))

	return httpServer, grpcServer, lis
}

func runServers(httpServer *http.Server, grpcServer *grpc.Server, lis net.Listener, log *zap.Logger) {
	var wg sync.WaitGroup
	errChan := make(chan error, ServersNumber)
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	wg.Add(1)
	go func(errChan chan error) {
		defer wg.Done()
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- fmt.Errorf("HTTP server error: %w", err)
		}
	}(errChan)

	wg.Add(1)
	go func(errChan chan error) {
		defer wg.Done()
		if err := grpcServer.Serve(lis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			errChan <- fmt.Errorf("gRPC server error: %w", err)
		}
	}(errChan)

	select {
	case <-stopChan:
		log.Info("Received shutdown signal, starting graceful shutdown")
	case err := <-errChan:
		log.Error("Server error: " + err.Error())
	}

	shutdownServers(httpServer, grpcServer, log)
	wg.Wait()
	log.Info("Shutdown complete")
}

func shutdownServers(httpServer *http.Server, grpcServer *grpc.Server, log *zap.Logger) {
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), shutdownDuration)
	defer shutdownCancel()

	var wg sync.WaitGroup
	wg.Add(ServersNumber)

	go func(httpServer *http.Server) {
		defer wg.Done()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			log.Error("HTTP server shutdown failed: " + err.Error())
		} else {
			log.Info("HTTP server shutdown gracefully")
		}
	}(httpServer)

	go func(grpcServer *grpc.Server) {
		defer wg.Done()
		grpcServer.GracefulStop()
		log.Info("gRPC server shutdown gracefully")
	}(grpcServer)

	wg.Wait()
}
