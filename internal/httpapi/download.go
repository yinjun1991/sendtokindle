package httpapi

import (
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func (h *Handlers) DownloadBook(c *gin.Context) {
	name := c.Param("name")

	f, info, err := h.Store.Open(name)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	defer func() { _ = f.Close() }()

	escaped := url.PathEscape(info.Name())
	c.Header("Content-Disposition", "attachment; filename*=UTF-8''"+escaped)
	c.Header("Cache-Control", "no-store")

	http.ServeContent(c.Writer, c.Request, filepath.Base(info.Name()), info.ModTime(), f)
}

