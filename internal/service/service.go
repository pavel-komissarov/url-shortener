package service

import (
	"errors"
	"fmt"

	"go.uber.org/zap"

	"url-shortener/internal/storage/errs"
	"url-shortener/pkg/util/random"
)

const shortURLLength = 10

type Storage interface {
	Put(url, shortURL string) error
	Get(url string) (string, error)
}

type Shortener struct {
	Storage Storage
	Log     *zap.Logger
}

func NewShortener(storage Storage, log *zap.Logger) *Shortener {
	return &Shortener{Storage: storage, Log: log}
}

func (s *Shortener) Shorten(url string) (string, error) {
	s.Log.Info("Shorten URL", zap.String("url", url))

	shortURL, err := random.NewRandomString(shortURLLength)
	if err != nil {
		return "", err
	}

	err = s.Storage.Put(url, shortURL)
	if err != nil {
		if errors.Is(err, errs.ErrURLIsExist) {
			return "", fmt.Errorf("url already exists")
		}
		return "", err
	}

	return shortURL, nil
}

func (s *Shortener) Resolve(url string) (string, error) {
	s.Log.Info("Resolve URL", zap.String("url", url))

	originURL, err := s.Storage.Get(url)
	if err != nil {
		if errors.Is(err, errs.ErrURLIsNotExist) {
			return "", fmt.Errorf("url does not exist")
		}
		return "", err
	}

	return originURL, nil
}
