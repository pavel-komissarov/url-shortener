package server

import (
	"context"
	"errors"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"

	"url-shortener/internal/grpc/urlshortener"
)

type Service interface {
	Resolve(url string) (string, error)
	Shorten(url string) (string, error)
}

type GRPCServer struct {
	urlshortener.UnimplementedURLShortenerServer
	Service Service
	Log     *zap.Logger
}

func (s *GRPCServer) Shorten(_ context.Context, req *urlshortener.ShortenRequest) (*urlshortener.ShortenResponse, error) {
	s.Log.Info("Shorten request", zap.String("url", req.GetUrl()))

	if err := validator.New().Var(req.Url, "required,url"); err != nil {
		s.Log.Error("Validation failed", zap.Error(err))
		return nil, errors.New("invalid URL format")
	}

	shortURL, err := s.Service.Shorten(req.Url)
	if err != nil {
		return nil, err
	}

	return &urlshortener.ShortenResponse{ShortUrl: shortURL}, nil
}

func (s *GRPCServer) Resolve(_ context.Context, req *urlshortener.ResolveRequest) (*urlshortener.ResolveResponse, error) {
	s.Log.Info("Resolve request", zap.String("short-URL", req.ShortUrl))

	url, err := s.Service.Resolve(req.ShortUrl)
	if err != nil {
		return nil, err
	}

	return &urlshortener.ResolveResponse{OriginalUrl: url}, nil
}
