package handler

import (
	"io/fs"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/zex/zpanel/internal/model"
	webembed "github.com/zex/zpanel/internal/web"
)

func (h *Handler) mountPanel(parent gin.IRouter, entryPrefix string) {
	staticRoot, err := fs.Sub(webembed.StaticFS, "dist")
	if err != nil {
		return
	}

	h.registerAPI(parent.Group("/api/v1"))

	serve := func(c *gin.Context) {
		if c.Request.Method != http.MethodGet {
			c.JSON(http.StatusNotFound, model.Fail("not found", "NOT_FOUND"))
			return
		}
		serveSPA(c, staticRoot, entryPrefix)
	}
	parent.GET("/", serve)
	parent.GET("/*filepath", serve)
}

func serveSPA(c *gin.Context, root fs.FS, entryPrefix string) {
	path := c.Request.URL.Path
	if entryPrefix != "" {
		path = strings.TrimPrefix(path, entryPrefix)
	}
	path = strings.TrimPrefix(path, "/")
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
		if entryPrefix != "" && path == "index.html" {
			inject := `<script>window.__ZPANEL_ENTRY__="` + entryPrefix + `"</script>`
			html := string(data)
			if strings.Contains(html, "</head>") {
				html = strings.Replace(html, "</head>", inject+"</head>", 1)
			} else {
				html = inject + html
			}
			data = []byte(html)
		}
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
