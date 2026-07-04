package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zex/zpanel/internal/model"
)

func (h *Handler) LNMPStatus(c *gin.Context) {
	c.JSON(http.StatusOK, model.OK(h.lnmp.Status()))
}

func (h *Handler) LNMPInstall(c *gin.Context) {
	res, err := h.lnmp.Install()
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Fail(err.Error(), "LNMP_INSTALL_FAILED"))
		return
	}
	_ = h.store.WriteAudit("lnmp_install", "success", c.ClientIP())
	c.JSON(http.StatusOK, model.OK(res))
}

func (h *Handler) ListSites(c *gin.Context) {
	sites, err := h.sites.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Fail(err.Error(), "INTERNAL_ERROR"))
		return
	}
	if sites == nil {
		sites = []model.Site{}
	}
	c.JSON(http.StatusOK, model.OK(sites))
}

func (h *Handler) GetSite(c *gin.Context) {
	site, err := h.sites.Get(c.Param("id"))
	if err != nil || site == nil {
		c.JSON(http.StatusNotFound, model.Fail("site not found", "SITE_NOT_FOUND"))
		return
	}
	c.JSON(http.StatusOK, model.OK(site))
}

func (h *Handler) CreateSite(c *gin.Context) {
	var req model.CreateSiteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Fail("invalid request", "BAD_REQUEST"))
		return
	}
	site, err := h.sites.Create(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Fail(err.Error(), "SITE_CREATE_FAILED"))
		return
	}
	_ = h.store.WriteAudit("site_create", site.Name, c.ClientIP())
	c.JSON(http.StatusOK, model.OK(site))
}

func (h *Handler) DeleteSite(c *gin.Context) {
	id := c.Param("id")
	site, _ := h.sites.Get(id)
	if err := h.sites.Delete(id); err != nil {
		c.JSON(http.StatusBadRequest, model.Fail(err.Error(), "SITE_DELETE_FAILED"))
		return
	}
	name := id
	if site != nil {
		name = site.Name
	}
	_ = h.store.WriteAudit("site_delete", name, c.ClientIP())
	c.JSON(http.StatusOK, model.OK(gin.H{"deleted": true}))
}
