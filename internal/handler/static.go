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
			html := rewriteIndexForEntry(string(data), entryPrefix)
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

// rewriteIndexForEntry 安全入口下修正静态资源路径，避免 /assets 404 导致白屏
func rewriteIndexForEntry(html, entryPrefix string) string {
	base := entryPrefix + "/"
	// 绝对路径 /assets/ → 相对 assets/（配合 base href）
	html = strings.ReplaceAll(html, `src="/assets/`, `src="assets/`)
	html = strings.ReplaceAll(html, `href="/assets/`, `href="assets/`)
	html = strings.ReplaceAll(html, `src="./assets/`, `src="assets/`)
	html = strings.ReplaceAll(html, `href="./assets/`, `href="assets/`)
	inject := `<base href="` + base + `"><script>window.__ZPANEL_ENTRY__="` + entryPrefix + `"</script>`
	if strings.Contains(html, "</head>") {
		return strings.Replace(html, "</head>", inject+"</head>", 1)
	}
	return inject + html
}
