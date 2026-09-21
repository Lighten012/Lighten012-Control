package podman

import (
	"net/http"

	"github.com/lighten012/control/internal/modules"
	"github.com/lighten012/control/internal/runner"
)

// Module 实现 modules.Module 接口的 Podman 组件模块。
type Module struct {
	handler *handler
}

// 编译期保证 *Module 实现了 modules.Module 接口。
var _ modules.Module = (*Module)(nil)

// NewModule 创建 Podman 组件模块。
func NewModule(r *runner.Runner) *Module {
	return &Module{handler: &handler{svc: NewService(r)}}
}

// Name 模块名称。
func (m *Module) Name() string { return "podman" }

// RegisterRoutes 注册 Podman 组件的只读查询路由（字面量，约定见 modules.Module）。
func (m *Module) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /lighten012-api/podman/containers", m.handler.handleListContainers)
	mux.HandleFunc("GET /lighten012-api/podman/images", m.handler.handleListImages)
	mux.HandleFunc("GET /lighten012-api/podman/networks", m.handler.handleListNetworks)
	mux.HandleFunc("GET /lighten012-api/podman/volumes", m.handler.handleListVolumes)
	mux.HandleFunc("GET /lighten012-api/podman/info", m.handler.handleInfo)
}
