package nginx

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
)

// TemplateData Nginx 模板数据
type TemplateData struct {
	Name       string
	Domains    []string
	Root       string
	PHPVersion string
	GoPort     int
	Binary     string
}

// Render 渲染 Nginx 配置模板
func Render(tmplPath string, data TemplateData) (string, error) {
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return "", fmt.Errorf("parse template: %w", err)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("execute template: %w", err)
	}
	return buf.String(), nil
}

// WriteConfig 写入配置文件
func WriteConfig(path, content string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(content), 0o644)
}

// EnableSite 软链到 sites-enabled
func EnableSite(available, enabled, name string) error {
	src := filepath.Join(available, name+".conf")
	dst := filepath.Join(enabled, name+".conf")
	_ = os.Remove(dst)
	return os.Symlink(src, dst)
}

// Test 执行 nginx -t
func Test() error {
	return run("nginx", "-t")
}

// Reload 重载 nginx
func Reload() error {
	if err := Test(); err != nil {
		return fmt.Errorf("nginx test failed: %w", err)
	}
	return RunSystemctl("reload", "nginx")
}

// RunSystemctl 执行 systemctl（白名单封装）
func RunSystemctl(args ...string) error {
	return run("systemctl", args...)
}

func run(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %s: %w", name, string(out), err)
	}
	return nil
}
