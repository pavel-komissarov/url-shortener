package server

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"

	"url-shortener/internal/grpc/urlshortener"
	"url-shortener/internal/service"
	"url-shortener/internal/storage/memory"
)

const (
	originalURL = "https://example.com"
	shortedURL  = "example.com"
)

func TestGRPCServer_Shorten_Success(t *testing.T) {
	logger, _ := zap.NewProduction()
	storage := memory.NewStorageInMemory(logger)
	shortenerService := service.NewShortener(storage, logger)
	grpcServer := &GRPCServer{Service: shortenerService, Log: logger}

	req := &urlshortener.ShortenRequest{Url: originalURL}
	resp, err := grpcServer.Shorten(context.Background(), req)

	assert.NoError(t, err)
	assert.NotEmpty(t, resp.ShortUrl)
}

func TestGRPCServer_Shorten_InvalidURL(t *testing.T) {
	logger, _ := zap.NewProduction()
	storage := memory.NewStorageInMemory(logger)
	shortenerService := service.NewShortener(storage, logger)
	grpcServer := &GRPCServer{Service: shortenerService, Log: logger}

	req := &urlshortener.ShortenRequest{Url: "invalid-url"}
	resp, err := grpcServer.Shorten(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestGRPCServer_Resolve_Success(t *testing.T) {
	logger, _ := zap.NewProduction()
	storage := memory.NewStorageInMemory(logger)
	shortenerService := service.NewShortener(storage, logger)
	grpcServer := &GRPCServer{Service: shortenerService, Log: logger}

	shortURL, err := shortenerService.Shorten(originalURL)
	assert.NoError(t, err)

	req := &urlshortener.ResolveRequest{ShortUrl: shortURL}
	resp, err := grpcServer.Resolve(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, originalURL, resp.OriginalUrl)
}

func TestGRPCServer_Resolve_NotFound(t *testing.T) {
	logger, _ := zap.NewProduction()
	storage := memory.NewStorageInMemory(logger)
	shortenerService := service.NewShortener(storage, logger)
	grpcServer := &GRPCServer{Service: shortenerService, Log: logger}

	req := &urlshortener.ResolveRequest{ShortUrl: shortedURL}
	resp, err := grpcServer.Resolve(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)
}
