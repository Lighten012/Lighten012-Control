// Lighten012-Control 后端服务入口。
//
// 启动流程：构建命令执行器 → 注册功能模块 → 组装路由 → 启动 HTTP 服务。
package main

import (
	"log"
	"net/http"
	"time"

	"github.com/lighten012/control/internal/httpapi"
	"github.com/lighten012/control/internal/modules/podman"
	"github.com/lighten012/control/internal/runner"
	"github.com/lighten012/control/web"
)

// 服务监听地址与模块命令执行超时，当前为固定值。
const (
	listenAddr = ":80"
	cmdTimeout = 5 * time.Second
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmsgprefix)
	log.SetPrefix("[lighten012] ")

	executor := runner.New(cmdTimeout)

	podmanModule := podman.NewModule(executor)

	server := &http.Server{
		Addr:              listenAddr,
		Handler:           httpapi.New(web.HTTP(), podmanModule.RegisterRoutes),
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("服务已启动: http://localhost%s （Ctrl+C 退出）", listenAddr)
	log.Fatal(server.ListenAndServe())
}
