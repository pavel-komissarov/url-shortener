package httpserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"url-shortener/internal/config"
	"url-shortener/internal/http/handlers/resolve"
	"url-shortener/internal/http/handlers/shorten"
	"url-shortener/internal/http/middleware/mvlogger"
)

type Service interface {
	Resolve(url string) (string, error)
	Shorten(url string) (string, error)
}

func NewHTTPServer(cfg config.ServerConfig, service Service, log *zap.Logger) *http.Server {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(mvlogger.NewLoggerMiddleware(log))

	r.POST("/shorten", shorten.New(service, log))
	r.GET("/resolve", resolve.New(service, log))

	server := &http.Server{
		Addr:         cfg.HTTPPort,
		Handler:      r,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	return server
}
