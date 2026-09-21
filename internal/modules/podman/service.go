// Package podman Podman 组件模块：以只读方式查询容器、镜像、网络、
// 存储卷与运行环境信息。
//
// 所有命令均为服务端约定的白名单命令（硬编码在代码内，不拼接任何
// 外部输入），通过 runner 执行并把 JSON 输出映射为 DTO，
// podman 原始结构不会泄漏到模块层之外。
package podman

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/lighten012/control/internal/httpapi"
	"github.com/lighten012/control/internal/runner"
)

// Service Podman 组件的业务逻辑。
type Service struct {
	runner *runner.Runner
}

// NewService 创建 Podman 组件服务。
func NewService(r *runner.Runner) *Service {
	return &Service{runner: r}
}

// query 执行白名单内的 podman 查询命令（自动追加 --format json）
// 并把输出解析为 T。
func query[T any](ctx context.Context, r *runner.Runner, args ...string) (T, error) {
	var out T
	cmd := append(append([]string{}, args...), "--format", "json")
	res, err := r.Run(ctx, "podman", cmd...)
	if err != nil {
		return out, fmt.Errorf("执行 podman %s 失败: %w", strings.Join(args, " "), err)
	}
	if err := json.Unmarshal([]byte(res.Stdout), &out); err != nil {
		return out, fmt.Errorf("解析 podman %s 输出失败: %w", strings.Join(args, " "), err)
	}
	return out, nil
}

// mapAll 把原始条目切片逐个映射为 DTO 切片。
func mapAll[F, T any](raw []F, mapFn func(F) T) []T {
	list := make([]T, 0, len(raw))
	for _, item := range raw {
		list = append(list, mapFn(item))
	}
	return list
}

// ListContainers 列出全部容器（含已停止）。podman ps -a
func (s *Service) ListContainers(ctx context.Context) ([]Container, error) {
	raw, err := query[[]podmanContainer](ctx, s.runner, "ps", "-a")
	if err != nil {
		return nil, err
	}
	return mapAll(raw, mapContainer), nil
}

// ListImages 列出本地镜像。podman images
func (s *Service) ListImages(ctx context.Context) ([]Image, error) {
	raw, err := query[[]podmanImage](ctx, s.runner, "images")
	if err != nil {
		return nil, err
	}
	return mapAll(raw, mapImage), nil
}

// ListNetworks 列出网络。podman network ls
func (s *Service) ListNetworks(ctx context.Context) ([]Network, error) {
	raw, err := query[[]podmanNetwork](ctx, s.runner, "network", "ls")
	if err != nil {
		return nil, err
	}
	return mapAll(raw, mapNetwork), nil
}

// ListVolumes 列出存储卷。podman volume ls
func (s *Service) ListVolumes(ctx context.Context) ([]Volume, error) {
	raw, err := query[[]podmanVolume](ctx, s.runner, "volume", "ls")
	if err != nil {
		return nil, err
	}
	return mapAll(raw, mapVolume), nil
}

// Info 查询运行环境信息。podman info
func (s *Service) Info(ctx context.Context) (Info, error) {
	raw, err := query[podmanInfo](ctx, s.runner, "info")
	if err != nil {
		return Info{}, err
	}
	return mapInfo(raw), nil
}

// ---------- HTTP 处理逻辑 ----------

// handler Podman 组件的 HTTP 处理逻辑。
type handler struct {
	svc *Service
}

// handleListContainers 处理 GET /lighten012-api/podman/containers。
func (h *handler) handleListContainers(w http.ResponseWriter, r *http.Request) {
	list, err := h.svc.ListContainers(r.Context())
	if err != nil {
		httpapi.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httpapi.WriteOK(w, list)
}

// handleListImages 处理 GET /lighten012-api/podman/images。
func (h *handler) handleListImages(w http.ResponseWriter, r *http.Request) {
	list, err := h.svc.ListImages(r.Context())
	if err != nil {
		httpapi.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httpapi.WriteOK(w, list)
}

// handleListNetworks 处理 GET /lighten012-api/podman/networks。
func (h *handler) handleListNetworks(w http.ResponseWriter, r *http.Request) {
	list, err := h.svc.ListNetworks(r.Context())
	if err != nil {
		httpapi.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httpapi.WriteOK(w, list)
}

// handleListVolumes 处理 GET /lighten012-api/podman/volumes。
func (h *handler) handleListVolumes(w http.ResponseWriter, r *http.Request) {
	list, err := h.svc.ListVolumes(r.Context())
	if err != nil {
		httpapi.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httpapi.WriteOK(w, list)
}

// handleInfo 处理 GET /lighten012-api/podman/info。
func (h *handler) handleInfo(w http.ResponseWriter, r *http.Request) {
	info, err := h.svc.Info(r.Context())
	if err != nil {
		httpapi.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httpapi.WriteOK(w, info)
}
