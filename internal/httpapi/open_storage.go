package httpapi

import (
	"net"
	"net/http"
	"os/exec"
	"runtime"

	"github.com/gin-gonic/gin"
)

func (h *Handlers) OpenStorageDir(c *gin.Context) {
	ip := net.ParseIP(c.ClientIP())
	if ip == nil || !ip.IsLoopback() {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	if err := openPath(h.StoreRoot); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func openPath(path string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", path)
	case "windows":
		cmd = exec.Command("explorer", path)
	default:
		cmd = exec.Command("xdg-open", path)
	}
	return cmd.Start()
}
