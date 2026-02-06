package httpapi

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"sendtokindle/internal/storage"
)

func TestBooks_CRUD(t *testing.T) {
	t.Parallel()

	store, err := storage.New(t.TempDir())
	if err != nil {
		t.Fatalf("init store: %v", err)
	}
	router := mustNewRouter(t, store)

	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	part, err := w.CreateFormFile("file", "a.epub")
	if err != nil {
		t.Fatalf("create part: %v", err)
	}
	if _, err := part.Write([]byte("content")); err != nil {
		t.Fatalf("write part: %v", err)
	}
	if err := w.Close(); err != nil {
		t.Fatalf("close writer: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/books", &body)
	req.Header.Set("Content-Type", w.FormDataContentType())
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("upload status = %d, want %d, body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}

	var uploadRes struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &uploadRes); err != nil {
		t.Fatalf("unmarshal upload response: %v", err)
	}
	if uploadRes.Name == "" {
		t.Fatalf("expected uploaded name")
	}

	req = httptest.NewRequest(http.MethodGet, "/api/books", nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("list status = %d, want %d", rec.Code, http.StatusOK)
	}
	if !bytes.Contains(rec.Body.Bytes(), []byte(uploadRes.Name)) {
		t.Fatalf("expected list response contains uploaded name, body=%s", rec.Body.String())
	}

	delBody, _ := json.Marshal(map[string]string{"name": uploadRes.Name})
	req = httptest.NewRequest(http.MethodPost, "/api/books/delete", bytes.NewReader(delBody))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("delete status = %d, want %d, body=%s", rec.Code, http.StatusNoContent, rec.Body.String())
	}
}

