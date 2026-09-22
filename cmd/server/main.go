// Lighten012-Control 后端服务入口。
//
// 启动流程：连接 Podman socket → 挂载模块路由 → 启动 HTTP 服务。
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/containers/podman/v5/pkg/bindings"
	"github.com/lighten012/control/internal/httpapi"
	"github.com/lighten012/control/internal/modules/podman"
	"github.com/lighten012/control/web"
)

// 服务监听地址与 Podman 请求超时，当前为固定值。
const (
	listenAddr     = ":80"
	requestTimeout = 5 * time.Second
)

// podmanURI 返回本机 Podman API socket 的连接地址。
// 优先使用 PODMAN_URI，rootless 默认取用户运行时目录。
func podmanURI() string {
	if uri := os.Getenv("PODMAN_URI"); uri != "" {
		return uri
	}
	if runtimeDir := os.Getenv("XDG_RUNTIME_DIR"); runtimeDir != "" {
		return "unix://" + filepath.Join(runtimeDir, "podman", "podman.sock")
	}
	if os.Getuid() != 0 {
		return "unix:///run/user/" + strconv.Itoa(os.Getuid()) + "/podman/podman.sock"
	}
	return "unix:///run/podman/podman.sock"
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lmsgprefix)
	log.SetPrefix("[lighten012] ")

	mux := http.NewServeMux()

	uri := podmanURI()
	_, err := bindings.NewConnection(context.Background(), uri)
	if err != nil {
		log.Fatalf("连接 Podman socket 失败: %v", err)
	}

	podman.NewModule(uri, requestTimeout).RegisterRoutes(mux)
	httpapi.MountAPI(mux)
	httpapi.MountStatic(mux, web.HTTP())

	server := &http.Server{
		Addr:              listenAddr,
		Handler:           httpapi.WithLogging(mux),
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("服务已启动: http://localhost%s （Ctrl+C 退出）", listenAddr)
	log.Fatal(server.ListenAndServe())
}
