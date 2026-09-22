package podman

import (
	"context"
	"net/http"

	"github.com/lighten012/control/internal/httpapi"
	"github.com/lighten012/control/internal/runner"
)

// Module 是 Podman 功能模块。
type Module struct {
	svc *Service
}

// NewModule 创建 Podman 组件模块。
func NewModule(r *runner.Executor) *Module {
	return &Module{svc: NewService(r)}
}

// RegisterRoutes 注册 Podman 组件的只读查询路由。
func (m *Module) RegisterRoutes(mux *http.ServeMux) {
	registerGet(mux, m.route("containers"), m.svc.ListContainers)
	registerGet(mux, m.route("images"), m.svc.ListImages)
	registerGet(mux, m.route("networks"), m.svc.ListNetworks)
	registerGet(mux, m.route("volumes"), m.svc.ListVolumes)
	registerGet(mux, m.route("info"), m.svc.Info)
}

// route 拼出 Podman 资源的完整 API 路径。
func (m *Module) route(resource string) string {
	return httpapi.PathPrefix + "/podman/" + resource
}

// registerGet 注册一个只读查询路由：执行查询并把结果按统一 JSON 格式写出。
func registerGet[T any](mux *http.ServeMux, path string, fn func(ctx context.Context) (T, error)) {
	mux.HandleFunc("GET "+path, func(w http.ResponseWriter, r *http.Request) {
		data, err := fn(r.Context())
		if err != nil {
			httpapi.WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}
		httpapi.WriteOK(w, data)
	})
}
