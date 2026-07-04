package app

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/zex/zpanel/internal/config"
	"github.com/zex/zpanel/internal/handler"
)

// Options 服务启动选项
type Options struct {
	ConfigPath string
	Version    string
}

// RunServer 启动 HTTP 面板服务
func RunServer(opts Options) error {
	configPath := opts.ConfigPath
	if configPath == "" {
		configPath = config.ResolvePath()
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		return err
	}
	if err := cfg.EnsureDirs(); err != nil {
		return err
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery(), gin.Logger())

	h := handler.New(cfg, opts.Version)
	h.Register(r)

	addr := cfg.ListenAddr()
	log.Printf("zpanel %s listening on http://%s", opts.Version, addr)
	log.Printf("config: %s", configPath)

	if err := r.Run(addr); err != nil {
		return fmt.Errorf("server: %w", err)
	}
	return nil
}
