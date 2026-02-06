package httpapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handlers) IndexPage(c *gin.Context) {
	books, err := h.Store.List()
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Header("Content-Type", "text/html; charset=utf-8")
	_ = h.Renderer.RenderIndex(c.Writer, struct {
		Books any
	}{Books: books})
}

func (h *Handlers) AdminPage(c *gin.Context) {
	c.Header("Content-Type", "text/html; charset=utf-8")
	_ = h.Renderer.RenderAdmin(c.Writer, struct {
		KindleURL string
		StoreRoot string
	}{
		KindleURL: h.KindleURL,
		StoreRoot: h.StoreRoot,
	})
}
