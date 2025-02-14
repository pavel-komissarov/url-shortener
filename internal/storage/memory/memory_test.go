package memory

import (
	"errors"
	"fmt"
	"sync"
	"testing"

	"go.uber.org/zap/zaptest"

	"url-shortener/internal/storage/errs"
)

const (
	originalURL = "https://example.com"
	shortedURL  = "exmpl"
)

func TestStorageInMemory_PutAndGet(t *testing.T) {
	t.Parallel()

	logger := zaptest.NewLogger(t)
	storage := NewStorageInMemory(logger)

	url := originalURL
	shortURL := shortedURL

	err := storage.Put(url, shortURL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	gotURL, err := storage.Get(shortURL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if gotURL != url {
		t.Errorf("got %v, want %v", gotURL, url)
	}
}

func TestStorageInMemory_PutDuplicate(t *testing.T) {
	t.Parallel()

	logger := zaptest.NewLogger(t)
	storage := NewStorageInMemory(logger)

	url := originalURL
	shortURL := shortedURL

	err := storage.Put(url, shortURL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = storage.Put(url, shortURL)
	if !errors.Is(err, errs.ErrURLIsExist) {
		t.Errorf("expected error %v, got %v", errs.ErrURLIsExist, err)
	}
}

func TestStorageInMemory_GetNotFound(t *testing.T) {
	t.Parallel()

	logger := zaptest.NewLogger(t)
	storage := NewStorageInMemory(logger)

	_, err := storage.Get("nonexistent")
	if !errors.Is(err, errs.ErrURLIsNotExist) {
		t.Errorf("expected error %v, got %v", errs.ErrURLIsNotExist, err)
	}
}

func TestStorageInMemory_ConcurrencyStress(t *testing.T) {
	t.Parallel()

	logger := zaptest.NewLogger(t)
	storage := NewStorageInMemory(logger)

	var wg sync.WaitGroup
	concurrentCount := 1000

	for i := 0; i < concurrentCount; i++ {
		wg.Add(2)

		go func(i int) {
			defer wg.Done()
			url := fmt.Sprintf("https://example.com/%d", i)
			shortURL := fmt.Sprintf("short%d", i)
			_ = storage.Put(url, shortURL)
		}(i)

		go func(i int) {
			defer wg.Done()
			shortURL := fmt.Sprintf("short%d", i)
			_, _ = storage.Get(shortURL)
		}(i)
	}

	wg.Wait()
}
