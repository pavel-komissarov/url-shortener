package shorten

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"url-shortener/internal/service"
	"url-shortener/internal/storage/memory"
)

func TestShortenHandler_Success(t *testing.T) {
	logger, _ := zap.NewProduction()

	storage := memory.NewStorageInMemory(logger)
	shortener := service.NewShortener(storage, logger)

	handler := New(shortener, logger)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodPost, "/shorten", strings.NewReader(`{"url": "https://example.com"}`))

	handler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "short_url")
}

func TestShortenHandler_InvalidRequest(t *testing.T) {
	logger, _ := zap.NewProduction()

	storage := memory.NewStorageInMemory(logger)
	shortener := service.NewShortener(storage, logger)

	handler := New(shortener, logger)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodPost, "/shorten", nil)

	handler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestShortenHandler_InvalidURLFormat(t *testing.T) {
	logger, _ := zap.NewProduction()

	storage := memory.NewStorageInMemory(logger)
	shortener := service.NewShortener(storage, logger)

	handler := New(shortener, logger)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodPost, "/shorten", nil)

	invalidRequestBody := `{"url": "not_a_url"}`
	c.Request.Body = io.NopCloser(strings.NewReader(invalidRequestBody))

	handler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
