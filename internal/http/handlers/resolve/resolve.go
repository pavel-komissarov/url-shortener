package resolve

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"net/http"
)

type Request struct {
	ShortenedURL string `json:"short_url" validate:"required"`
}

type Response struct {
	URL    string `json:"original_url,omitempty"`
	Error  string `json:"error,omitempty"`
	Status string `json:"status"`
}

type Resolver interface {
	Resolve(url string) (string, error)
}

func New(service Resolver, log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		log := log.With(zap.String("op", "resolve"))

		var req Request
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Error("invalid request", zap.Error(err))
			c.JSON(http.StatusBadRequest, Response{Error: "invalid request", Status: "Error"})
			return
		}

		if err := validator.New().Struct(req); err != nil {
			log.Error("validation failed", zap.Error(err))
			c.JSON(http.StatusBadRequest, Response{Error: "invalid shortened URL", Status: "Error"})
			return
		}

		url, err := service.Resolve(req.ShortenedURL)
		if err != nil {
			log.Error("failed to resolve URL", zap.Error(err))
			c.JSON(http.StatusNotFound, Response{Error: "URL not found", Status: "Error"})
			return
		}

		c.JSON(http.StatusOK, Response{URL: url, Status: "OK"})
	}
}
