package store

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/zex/zpanel/internal/model"
)

func (s *Store) migrateSites() error {
	_, err := s.db.Exec(`
CREATE TABLE IF NOT EXISTS sites (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    type TEXT NOT NULL,
    status TEXT NOT NULL,
    domains TEXT NOT NULL,
    root TEXT NOT NULL,
    php_version TEXT,
    go_port INTEGER,
    go_binary TEXT,
    systemd_unit TEXT,
    nginx_config_path TEXT NOT NULL,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);
`)
	return err
}

// ListSites 列出所有站点
func (s *Store) ListSites() ([]model.Site, error) {
	rows, err := s.db.Query(`SELECT id, name, type, status, domains, root, php_version, go_port, go_binary, systemd_unit, nginx_config_path, created_at, updated_at FROM sites ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sites []model.Site
	for rows.Next() {
		site, err := scanSite(rows)
		if err != nil {
			return nil, err
		}
		sites = append(sites, *site)
	}
	return sites, rows.Err()
}

// GetSite 按 ID 查询
func (s *Store) GetSite(id string) (*model.Site, error) {
	row := s.db.QueryRow(`SELECT id, name, type, status, domains, root, php_version, go_port, go_binary, systemd_unit, nginx_config_path, created_at, updated_at FROM sites WHERE id = ?`, id)
	return scanSiteRow(row)
}

func scanSite(rows *sql.Rows) (*model.Site, error) {
	var site model.Site
	var domainsJSON, phpVer, goBin, systemd, created, updated sql.NullString
	var goPort sql.NullInt64
	if err := rows.Scan(&site.ID, &site.Name, &site.Type, &site.Status, &domainsJSON, &site.Root, &phpVer, &goPort, &goBin, &systemd, &site.NginxConfigPath, &created, &updated); err != nil {
		return nil, err
	}
	fillSite(&site, domainsJSON, phpVer, goPort, goBin, systemd, created, updated)
	return &site, nil
}

func scanSiteRow(row *sql.Row) (*model.Site, error) {
	var site model.Site
	var domainsJSON, phpVer, goBin, systemd, created, updated sql.NullString
	var goPort sql.NullInt64
	if err := row.Scan(&site.ID, &site.Name, &site.Type, &site.Status, &domainsJSON, &site.Root, &phpVer, &goPort, &goBin, &systemd, &site.NginxConfigPath, &created, &updated); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	fillSite(&site, domainsJSON, phpVer, goPort, goBin, systemd, created, updated)
	return &site, nil
}

func fillSite(site *model.Site, domainsJSON, phpVer sql.NullString, goPort sql.NullInt64, goBin, systemd, created, updated sql.NullString) {
	_ = json.Unmarshal([]byte(domainsJSON.String), &site.Domains)
	site.PHPVersion = phpVer.String
	if goPort.Valid {
		site.GoPort = int(goPort.Int64)
	}
	site.GoBinary = goBin.String
	site.SystemdUnit = systemd.String
	site.CreatedAt, _ = time.Parse(time.RFC3339, created.String)
	site.UpdatedAt, _ = time.Parse(time.RFC3339, updated.String)
}

// InsertSite 写入站点
func (s *Store) InsertSite(site *model.Site) error {
	domainsJSON, _ := json.Marshal(site.Domains)
	now := time.Now().UTC().Format(time.RFC3339)
	if site.CreatedAt.IsZero() {
		site.CreatedAt, _ = time.Parse(time.RFC3339, now)
	}
	site.UpdatedAt, _ = time.Parse(time.RFC3339, now)
	_, err := s.db.Exec(`INSERT INTO sites (id, name, type, status, domains, root, php_version, go_port, go_binary, systemd_unit, nginx_config_path, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		site.ID, site.Name, site.Type, site.Status, string(domainsJSON), site.Root, nullStr(site.PHPVersion), nullInt(site.GoPort), nullStr(site.GoBinary), nullStr(site.SystemdUnit), site.NginxConfigPath, site.CreatedAt.Format(time.RFC3339), site.UpdatedAt.Format(time.RFC3339))
	return err
}

// DeleteSite 删除站点记录
func (s *Store) DeleteSite(id string) error {
	_, err := s.db.Exec(`DELETE FROM sites WHERE id = ?`, id)
	return err
}

func nullStr(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}

func nullInt(n int) interface{} {
	if n == 0 {
		return nil
	}
	return n
}

// UpdateSiteStatus 更新状态
func (s *Store) UpdateSiteStatus(id string, status model.SiteStatus) error {
	now := time.Now().UTC().Format(time.RFC3339)
	res, err := s.db.Exec(`UPDATE sites SET status = ?, updated_at = ? WHERE id = ?`, status, now, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("site not found")
	}
	return nil
}
