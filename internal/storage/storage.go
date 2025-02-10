package storage

import (
	"go.uber.org/zap"

	"url-shortener/internal/config"
	"url-shortener/internal/storage/memory"
	"url-shortener/internal/storage/postgres"
)

type Storage interface {
	Put(url, shortURL string) error
	Get(url string) (string, error)
}

func NewStorage(storageConf *config.StorageConfig, log *zap.Logger) (Storage, error) {
	switch storageConf.Type {
	case "postgres":
		return postgres.NewStorage(storageConf.Postgres, log)
	default:
		return memory.NewStorageInMemory(log), nil
	}
}
