//go:build !linux

package site

import (
	"fmt"

	"github.com/zex/zpanel/internal/config"
	"github.com/zex/zpanel/internal/model"
	"github.com/zex/zpanel/internal/store"
)

// Service 站点管理（非 Linux 平台仅支持列表）
type Service struct {
	cfg   *config.Config
	store *store.Store
}

func NewService(cfg *config.Config, st *store.Store) *Service {
	return &Service{cfg: cfg, store: st}
}

func (s *Service) List() ([]model.Site, error) {
	return s.store.ListSites()
}

func (s *Service) Get(id string) (*model.Site, error) {
	return s.store.GetSite(id)
}

func (s *Service) Create(_ model.CreateSiteRequest) (*model.Site, error) {
	return nil, fmt.Errorf("site management requires Linux")
}

func (s *Service) Delete(_ string) error {
	return fmt.Errorf("site management requires Linux")
}
