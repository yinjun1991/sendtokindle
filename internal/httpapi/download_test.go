package httpapi

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"sendtokindle/internal/storage"
)

func TestDownloadBook_Range(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	store, err := storage.New(dir)
	if err != nil {
		t.Fatalf("init store: %v", err)
	}

	const filename = "hello.txt"
	content := []byte("hello world")
	if err := os.WriteFile(filepath.Join(dir, filename), content, 0o600); err != nil {
		t.Fatalf("write file: %v", err)
	}

	router := mustNewRouter(t, store)

	req := httptest.NewRequest(http.MethodGet, "/books/"+filename, nil)
	req.Header.Set("Range", "bytes=0-4")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusPartialContent {
		t.Fatalf("status = %d, want %d", res.StatusCode, http.StatusPartialContent)
	}
	body := w.Body.Bytes()
	if !bytes.Equal(body, []byte("hello")) {
		t.Fatalf("body = %q, want %q", string(body), "hello")
	}
	if res.Header.Get("Accept-Ranges") == "" {
		t.Fatalf("expected Accept-Ranges header")
	}
}

