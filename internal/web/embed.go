package web

import "embed"

// StaticFS 嵌入前端构建产物（开发期占位；Vue build 输出复制到 internal/web/dist）
//
//go:embed all:dist
var StaticFS embed.FS
