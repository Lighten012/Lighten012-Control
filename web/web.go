// Package web 托管内置的前端静态页面，随二进制一同编译，
// 保证产物单文件可随处运行。
package web

import (
	"embed"
	"net/http"
)

//go:embed index.html style.css
var files embed.FS

// HTTP 返回静态页面资源。
func HTTP() http.FileSystem {
	return http.FS(files)
}
