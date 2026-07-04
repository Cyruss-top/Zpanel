package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zex/zpanel/internal/auth"
	"github.com/zex/zpanel/internal/config"
	"github.com/zex/zpanel/internal/service/lnmp"
	"github.com/zex/zpanel/internal/service/site"
	"github.com/zex/zpanel/internal/store"
)

// Handler HTTP 路由处理器
type Handler struct {
	cfg     *config.Config
	version string
	store   *store.Store
	auth    *auth.Service
	lnmp    *lnmp.Service
	sites   *site.Service
}

func New(cfg *config.Config, version string, st *store.Store, authSvc *auth.Service) *Handler {
	return &Handler{
		cfg:     cfg,
		version: version,
		store:   st,
		auth:    authSvc,
		lnmp:    lnmp.NewService(),
		sites:   site.NewService(cfg, st),
	}
}

func (h *Handler) registerAPI(r gin.IRouter) {
	r.GET("/health", h.Health)
	r.POST("/auth/login", h.Login)

	authed := r.Group("")
	authed.Use(h.auth.Middleware())
	authed.GET("/monitor/overview", h.MonitorOverview)
	authed.GET("/lnmp/status", h.LNMPStatus)
	authed.POST("/lnmp/install", h.LNMPInstall)
	authed.GET("/sites", h.ListSites)
	authed.POST("/sites", h.CreateSite)
	authed.GET("/sites/:id", h.GetSite)
	authed.DELETE("/sites/:id", h.DeleteSite)
}

// Register 注册所有路由
func (h *Handler) Register(r *gin.Engine) {
	entry := h.cfg.EntryPrefix()

	if entry != "" {
		r.GET("/", func(c *gin.Context) {
			c.Redirect(http.StatusFound, entry+"/")
		})
		r.GET(entry, func(c *gin.Context) {
			c.Redirect(http.StatusFound, entry+"/")
		})
		h.mountPanel(r.Group(entry), entry)
	} else {
		h.mountPanel(r, "")
	}
	h.registerStaticFallback(r, entry)
}
