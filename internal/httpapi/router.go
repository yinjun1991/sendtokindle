package httpapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Config struct {
	GinMode string
}

func NewRouter(cfg Config, handlers *Handlers) *gin.Engine {
	if cfg.GinMode != "" {
		gin.SetMode(cfg.GinMode)
	}

	r := gin.New()
	r.Use(gin.Recovery())

	r.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	r.GET("/", handlers.IndexPage)
	r.GET("/admin", handlers.AdminPage)

	api := r.Group("/api")
	{
		api.GET("/books", handlers.ListBooks)
		api.POST("/books", handlers.UploadBook)
		api.POST("/books/delete", handlers.DeleteBook)
		api.POST("/open-storage", handlers.OpenStorageDir)
	}

	r.GET("/books/:name", handlers.DownloadBook)

	return r
}
