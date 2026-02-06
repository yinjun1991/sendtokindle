package httpapi

import (
	"testing"

	"github.com/gin-gonic/gin"

	"sendtokindle/internal/storage"
	"sendtokindle/internal/web"
)

func mustNewRouter(t *testing.T, store *storage.Store) *gin.Engine {
	t.Helper()

	renderer, err := web.NewRenderer()
	if err != nil {
		t.Fatalf("init renderer: %v", err)
	}

	handlers := &Handlers{
		Store:    store,
		Renderer: renderer,
	}
	return NewRouter(Config{GinMode: gin.TestMode}, handlers)
}
