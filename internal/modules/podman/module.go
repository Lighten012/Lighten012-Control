package podman

import (
	"net/http"
	"time"

	"github.com/lighten012/control/internal/httpapi"
)

// Module Podman 组件模块。
type Module struct {
	svc *Service
}

// NewModule 创建 Podman 组件模块。
func NewModule(uri string, timeout time.Duration) *Module {
	return &Module{svc: NewService(uri, timeout)}
}

// RegisterRoutes 注册 Podman 组件的只读查询路由。
func (m *Module) RegisterRoutes(mux *http.ServeMux) {
	base := httpapi.PathPrefix + "/podman"
	httpapi.RegisterGet(mux, base+"/containers", m.svc.ListContainers)
	httpapi.RegisterGet(mux, base+"/images", m.svc.ListImages)
	httpapi.RegisterGet(mux, base+"/networks", m.svc.ListNetworks)
	httpapi.RegisterGet(mux, base+"/volumes", m.svc.ListVolumes)
	httpapi.RegisterGet(mux, base+"/info", m.svc.Info)
}
