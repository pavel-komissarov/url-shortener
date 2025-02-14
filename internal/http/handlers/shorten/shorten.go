package shorten

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type Request struct {
	URL string `json:"url" validate:"required,url"`
}

type Response struct {
	ShortenedURL string `json:"short_url,omitempty"`
	Error        string `json:"error,omitempty"`
	Status       string `json:"status"`
}

type Shortener interface {
	Shorten(url string) (string, error)
}

func New(service Shortener, log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		log := log.With(zap.String("op", "shorten"))

		log.Debug("Shorten called")

		var req Request
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Error("invalid request", zap.Error(err))
			c.JSON(http.StatusBadRequest, Response{Error: "invalid request", Status: "Error"})
			return
		}

		log.Info("shorten request", zap.String("url", req.URL))

		if err := validator.New().Struct(req); err != nil {
			log.Error("validation failed", zap.Error(err))
			c.JSON(http.StatusBadRequest, Response{Error: "invalid URL format", Status: "Error"})
			return
		}

		shortened, err := service.Shorten(req.URL)
		if err != nil {
			log.Error("failed to shorten URL", zap.Error(err))
			c.JSON(http.StatusInternalServerError, Response{Error: err.Error(), Status: "Error"})
			return
		}

		c.JSON(http.StatusOK, Response{ShortenedURL: shortened, Status: "OK"})
	}
}
