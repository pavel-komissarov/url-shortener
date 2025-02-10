package grpcserver

import (
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"url-shortener/internal/config"
	grpcShortoner "url-shortener/internal/grpc/server"
	"url-shortener/internal/grpc/urlshortener"
)

type Service interface {
	Resolve(url string) (string, error)
	Shorten(url string) (string, error)
}

func New(cfg config.ServerConfig, service Service, log *zap.Logger) *grpc.Server {
	server := grpc.NewServer(
		grpc.ConnectionTimeout(cfg.Timeout),
	)

	grpcServer := &grpcShortoner.GRPCServer{Service: service, Log: log}
	urlshortener.RegisterURLShortenerServer(server, grpcServer)

	return server
}
