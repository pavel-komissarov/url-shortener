package memory

import (
	"sync"

	"go.uber.org/zap"

	"url-shortener/internal/storage/errs"
)

type StorageInMemory struct {
	rvMu    sync.RWMutex
	storage map[string]string
	reverse map[string]string
	log     *zap.Logger
}

func NewStorageInMemory(log *zap.Logger) *StorageInMemory {
	return &StorageInMemory{
		storage: make(map[string]string),
		reverse: make(map[string]string),
		log:     log,
	}
}

func (s *StorageInMemory) Put(url, shortURL string) error {
	s.rvMu.Lock()
	defer s.rvMu.Unlock()

	s.log.Debug("put", zap.String("url", url), zap.String("shortUrl", shortURL))

	if _, ok := s.storage[shortURL]; ok {
		return errs.ErrURLIsExist
	}

	if _, ok := s.reverse[url]; ok {
		return errs.ErrURLIsExist
	}

	s.storage[shortURL] = url
	s.reverse[url] = shortURL

	return nil
}

func (s *StorageInMemory) Get(shortURL string) (string, error) {
	s.rvMu.RLock()
	defer s.rvMu.RUnlock()

	s.log.Debug("get", zap.String("shortUrl", shortURL))

	if url, ok := s.storage[shortURL]; ok {
		return url, nil
	}

	return "", errs.ErrURLIsNotExist
}
