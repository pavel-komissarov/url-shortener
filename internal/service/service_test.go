package service

import (
	"go.uber.org/zap"
	"testing"

	"github.com/stretchr/testify/assert"
	"url-shortener/internal/storage/memory"
)

const (
	originalURL = "https://example.com"
)

func TestShorten(t *testing.T) {
	logger, _ := zap.NewProduction()

	storage := memory.NewStorageInMemory(logger)
	service := NewShortener(storage, logger)

	shortURL, err := service.Shorten(originalURL)
	assert.NoError(t, err)
	assert.Len(t, shortURL, 10)

	resolvedURL, err := storage.Get(shortURL)
	assert.NoError(t, err)
	assert.Equal(t, originalURL, resolvedURL)
}

func TestResolve(t *testing.T) {
	logger, _ := zap.NewProduction()

	storage := memory.NewStorageInMemory(logger)
	service := NewShortener(storage, logger)

	shortURL, err := service.Shorten(originalURL)
	assert.NoError(t, err)

	url, err := service.Resolve(shortURL)
	assert.NoError(t, err)
	assert.Equal(t, originalURL, url)
}

func TestShorten_UrlExists(t *testing.T) {
	logger, _ := zap.NewProduction()

	storage := memory.NewStorageInMemory(logger)
	service := NewShortener(storage, logger)

	_, err := service.Shorten(originalURL)
	assert.NoError(t, err)

	_, err = service.Shorten(originalURL)
	assert.Error(t, err)
	assert.Equal(t, "url already exists", err.Error())
}

func TestResolve_UrlNotExist(t *testing.T) {
	logger, _ := zap.NewProduction()

	storage := memory.NewStorageInMemory(logger)
	service := NewShortener(storage, logger)

	_, err := service.Resolve("nonexistent")
	assert.Error(t, err)
	assert.Equal(t, "url does not exist", err.Error())
}
