package httpapi

import (
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

const maxUploadBytes int64 = 512 << 20 // 512 MiB

func (h *Handlers) ListBooks(c *gin.Context) {
	books, err := h.Store.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "list books failed"})
		return
	}
	c.JSON(http.StatusOK, books)
}

func (h *Handlers) UploadBook(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxUploadBytes)

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "file too large"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing file"})
		return
	}
	defer func() { _ = file.Close() }()

	limited := &io.LimitedReader{R: file, N: maxUploadBytes + 1}
	savedName, err := h.Store.Save(header.Filename, limited)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "save file failed"})
		return
	}
	if limited.N == 0 {
		_ = h.Store.Delete(savedName)
		c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "file too large"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"name": savedName})
}

type deleteBookRequest struct {
	Name string `json:"name"`
}

func (h *Handlers) DeleteBook(c *gin.Context) {
	var req deleteBookRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if err := h.Store.Delete(req.Name); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.Status(http.StatusNoContent)
}
