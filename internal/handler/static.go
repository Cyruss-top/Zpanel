package handler

import (
	"io/fs"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/zex/zpanel/internal/model"
	webembed "github.com/zex/zpanel/internal/web"
)

func (h *Handler) registerStatic(r *gin.Engine) {
	staticRoot, err := fs.Sub(webembed.StaticFS, "dist")
	if err != nil {
		return
	}

	r.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api/") {
			c.JSON(http.StatusNotFound, model.Fail("not found", "NOT_FOUND"))
			return
		}
		if c.Request.Method != http.MethodGet {
			c.JSON(http.StatusNotFound, model.Fail("not found", "NOT_FOUND"))
			return
		}
		serveSPA(c, staticRoot)
	})
}

func serveSPA(c *gin.Context, root fs.FS) {
	path := strings.TrimPrefix(c.Request.URL.Path, "/")
	if path == "" {
		path = "index.html"
	}
	if _, err := fs.Stat(root, path); err != nil {
		path = "index.html"
	}
	data, err := fs.ReadFile(root, path)
	if err != nil {
		c.JSON(http.StatusNotFound, model.Fail("not found", "NOT_FOUND"))
		return
	}
	ctype := "text/plain"
	switch {
	case strings.HasSuffix(path, ".html"):
		ctype = "text/html; charset=utf-8"
	case strings.HasSuffix(path, ".js"):
		ctype = "application/javascript"
	case strings.HasSuffix(path, ".css"):
		ctype = "text/css"
	case strings.HasSuffix(path, ".svg"):
		ctype = "image/svg+xml"
	case strings.HasSuffix(path, ".ico"):
		ctype = "image/x-icon"
	}
	c.Data(http.StatusOK, ctype, data)
}
