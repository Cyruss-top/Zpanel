package store

import (
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

// User 面板管理员
type User struct {
	ID           int64
	Username     string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Store SQLite 数据访问
type Store struct {
	db *sql.DB
}

// Open 打开数据库并迁移 schema
func Open(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}
	db.SetMaxOpenConns(1)

	s := &Store{db: db}
	if err := s.migrate(); err != nil {
		db.Close()
		return nil, err
	}
	return s, nil
}

func (s *Store) migrate() error {
	schema := `
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS audit_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    action TEXT NOT NULL,
    detail TEXT,
    ip TEXT,
    created_at TEXT NOT NULL
);
`
	if _, err := s.db.Exec(schema); err != nil {
		return fmt.Errorf("migrate: %w", err)
	}
	return s.migrateSites()
}

func (s *Store) Close() error {
	return s.db.Close()
}

// GetUserByUsername 按用户名查询
func (s *Store) GetUserByUsername(username string) (*User, error) {
	row := s.db.QueryRow(
		`SELECT id, username, password_hash, created_at, updated_at FROM users WHERE username = ?`,
		username,
	)
	var u User
	var created, updated string
	if err := row.Scan(&u.ID, &u.Username, &u.PasswordHash, &created, &updated); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	u.CreatedAt, _ = time.Parse(time.RFC3339, created)
	u.UpdatedAt, _ = time.Parse(time.RFC3339, updated)
	return &u, nil
}

// UpsertUser 创建或更新用户密码
func (s *Store) UpsertUser(username, passwordHash string) error {
	now := time.Now().UTC().Format(time.RFC3339)
	existing, err := s.GetUserByUsername(username)
	if err != nil {
		return err
	}
	if existing == nil {
		_, err = s.db.Exec(
			`INSERT INTO users (username, password_hash, created_at, updated_at) VALUES (?, ?, ?, ?)`,
			username, passwordHash, now, now,
		)
		return err
	}
	_, err = s.db.Exec(
		`UPDATE users SET password_hash = ?, updated_at = ? WHERE username = ?`,
		passwordHash, now, username,
	)
	return err
}

// WriteAudit 记录审计日志
func (s *Store) WriteAudit(action, detail, ip string) error {
	now := time.Now().UTC().Format(time.RFC3339)
	_, err := s.db.Exec(
		`INSERT INTO audit_logs (action, detail, ip, created_at) VALUES (?, ?, ?, ?)`,
		action, detail, ip, now,
	)
	return err
}
