package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zex/zpanel/internal/auth"
	"github.com/zex/zpanel/internal/config"
	"github.com/zex/zpanel/internal/model"
	"github.com/zex/zpanel/internal/service/monitor"
	"github.com/zex/zpanel/internal/store"
)

// Handler HTTP 路由处理器
type Handler struct {
	cfg     *config.Config
	version string
	store   *store.Store
	auth    *auth.Service
}

func New(cfg *config.Config, version string, st *store.Store, authSvc *auth.Service) *Handler {
	return &Handler{cfg: cfg, version: version, store: st, auth: authSvc}
}

// Register 注册所有路由
func (h *Handler) Register(r *gin.Engine) {
	api := r.Group("/api/v1")
	api.GET("/health", h.Health)
	api.POST("/auth/login", h.Login)

	authed := api.Group("")
	authed.Use(h.auth.Middleware())
	authed.GET("/monitor/overview", h.MonitorOverview)

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

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Fail("invalid request", "BAD_REQUEST"))
		return
	}

	key := c.ClientIP() + ":" + req.Username
	if err := h.auth.AllowLogin(key); err != nil {
		c.JSON(http.StatusTooManyRequests, model.Fail("too many attempts", "RATE_LIMITED"))
		return
	}

	user, err := h.store.GetUserByUsername(req.Username)
	if err != nil || user == nil || !auth.CheckPassword(user.PasswordHash, req.Password) {
		_ = h.store.WriteAudit("login_failed", req.Username, c.ClientIP())
		c.JSON(http.StatusUnauthorized, model.Fail("invalid username or password", "INVALID_CREDENTIALS"))
		return
	}

	h.auth.ResetLoginAttempts(key)
	token, err := h.auth.GenerateToken(user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Fail("token error", "INTERNAL_ERROR"))
		return
	}
	_ = h.store.WriteAudit("login_success", user.Username, c.ClientIP())

	c.JSON(http.StatusOK, model.OK(gin.H{
		"token":    token,
		"username": user.Username,
		"expires":  "24h",
	}))
}

func (h *Handler) MonitorOverview(c *gin.Context) {
	overview, err := monitor.Collect(h.cfg.Paths.Data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Fail(err.Error(), "MONITOR_ERROR"))
		return
	}
	c.JSON(http.StatusOK, model.OK(overview))
}
