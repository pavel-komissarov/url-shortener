package mvlogger

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

func NewLoggerMiddleware(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		logger := log.With(
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
		)

		logger.Info("incoming request")

		c.Next()

		logger.Info("response sent", zap.String("duration", time.Since(start).String()), zap.Int("status", c.Writer.Status()))
	}
}
