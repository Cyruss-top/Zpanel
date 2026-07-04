package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zex/zpanel/internal/config"
	"github.com/zex/zpanel/internal/model"
)

// Handler HTTP 路由处理器
type Handler struct {
	cfg     *config.Config
	version string
}

func New(cfg *config.Config, version string) *Handler {
	return &Handler{cfg: cfg, version: version}
}

// Register 注册所有路由
func (h *Handler) Register(r *gin.Engine) {
	r.GET("/api/v1/health", h.Health)

	// 前端静态资源（占位，v0.3.0 完善）
	r.NoRoute(func(c *gin.Context) {
		if c.Request.Method != http.MethodGet {
			c.JSON(http.StatusNotFound, model.Fail("not found", "NOT_FOUND"))
			return
		}
		c.JSON(http.StatusNotFound, model.Fail("not found", "NOT_FOUND"))
	})
}

func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, model.OK(gin.H{
		"status":  "ok",
		"version": h.version,
	}))
}
