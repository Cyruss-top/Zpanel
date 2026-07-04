package app

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/zex/zpanel/internal/auth"
	"github.com/zex/zpanel/internal/config"
	"github.com/zex/zpanel/internal/handler"
	"github.com/zex/zpanel/internal/store"
)

// Options 服务启动选项
type Options struct {
	ConfigPath string
	Version    string
}

const defaultAdminPassword = "admin"

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

	st, err := store.Open(cfg.SQLitePath())
	if err != nil {
		return err
	}
	defer st.Close()

	if err := bootstrapAdmin(st, cfg); err != nil {
		return err
	}

	authSvc, err := auth.NewService(cfg.Paths.Data)
	if err != nil {
		return err
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery(), gin.Logger())

	h := handler.New(cfg, opts.Version, st, authSvc)
	h.Register(r)

	addr := cfg.ListenAddr()
	log.Printf("zpanel %s listening on http://%s", opts.Version, addr)
	log.Printf("config: %s", configPath)

	if err := r.Run(addr); err != nil {
		return fmt.Errorf("server: %w", err)
	}
	return nil
}

func bootstrapAdmin(st *store.Store, cfg *config.Config) error {
	user, err := st.GetUserByUsername(cfg.Auth.Username)
	if err != nil {
		return err
	}
	if user != nil {
		return nil
	}

	hash := cfg.Auth.PasswordHash
	if hash == "" {
		hash, err = auth.HashPassword(defaultAdminPassword)
		if err != nil {
			return err
		}
		log.Printf("WARN: created default admin user %q with password %q — change immediately",
			cfg.Auth.Username, defaultAdminPassword)
	}
	return st.UpsertUser(cfg.Auth.Username, hash)
}
