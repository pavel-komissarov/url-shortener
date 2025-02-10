package resolve

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

const (
	originalURL = "https://example.com"
)

func TestResolveHandler_Success(t *testing.T) {
	logger, _ := zap.NewProduction()
	storage := memory.NewStorageInMemory(logger)
	shortener := service.NewShortener(storage, logger)

	shortURL, err := shortener.Shorten(originalURL)
	assert.NoError(t, err)

	handler := New(shortener, logger)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodPost, "/resolve", nil)
	c.Request.Body = io.NopCloser(strings.NewReader(`{"short_url": "` + shortURL + `"}`))

	handler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), originalURL)
}

func TestResolveHandler_InvalidRequest(t *testing.T) {
	logger, _ := zap.NewProduction()

	storage := memory.NewStorageInMemory(logger)
	shortener := service.NewShortener(storage, logger)

	handler := New(shortener, logger)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodPost, "/resolve", nil)

	handler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestResolveHandler_InvalidShortenedURL(t *testing.T) {
	logger, _ := zap.NewProduction()

	storage := memory.NewStorageInMemory(logger)
	shortener := service.NewShortener(storage, logger)

	handler := New(shortener, logger)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodPost, "/resolve", nil)

	invalidRequestBody := `{"short_url": "nonexistent"}`
	c.Request.Body = io.NopCloser(strings.NewReader(invalidRequestBody))

	handler(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestResolveHandler_ValidationError(t *testing.T) {
	logger, _ := zap.NewProduction()

	storage := memory.NewStorageInMemory(logger)
	shortener := service.NewShortener(storage, logger)

	handler := New(shortener, logger)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodPost, "/resolve", nil)

	c.Request.Body = io.NopCloser(strings.NewReader(`{"short_url": ""}`))

	handler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestResolveHandler_UrlNotFound(t *testing.T) {
	logger, _ := zap.NewProduction()

	storage := memory.NewStorageInMemory(logger)
	shortener := service.NewShortener(storage, logger)

	handler := New(shortener, logger)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodPost, "/resolve", nil)

	invalidRequestBody := `{"short_url": "nonexistent"}`
	c.Request.Body = io.NopCloser(strings.NewReader(invalidRequestBody))

	handler(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
