// Package modules 定义 Linux 功能模块的统一抽象。
//
// 每个功能模块（podman、iptables 等）实现 Module 接口，
// 在程序启动时集中注册；HTTP 路由层只面向该接口，不感知具体模块实现，
// 从而保证新增模块时无需改动核心路由与基础设施代码。
package modules

import "net/http"

// Module 一个可被后端托管的 Linux 功能模块。
type Module interface {
	// Name 模块唯一名称，用于日志与路由归属标识，例如 "podman"。
	Name() string
	// RegisterRoutes 将模块的 API 路由注册到 mux。
	// 路径约定（口头约定，双方确认即可）：统一为 /lighten012-api/<操控组件>/...，
	// 具体子路径由各模块以字面量自行定义。
	RegisterRoutes(mux *http.ServeMux)
}

// Registry 模块注册表。
type Registry struct {
	modules []Module
	byName  map[string]Module
}

// NewRegistry 创建空的模块注册表。
func NewRegistry() *Registry {
	return &Registry{byName: make(map[string]Module)}
}

// Register 注册模块；名称重复视为程序缺陷，启动期立即 panic 暴露。
func (rg *Registry) Register(m Module) {
	if _, dup := rg.byName[m.Name()]; dup {
		panic("modules: 重复注册模块 " + m.Name())
	}
	rg.byName[m.Name()] = m
	rg.modules = append(rg.modules, m)
}

// Modules 返回已注册的全部模块（按注册顺序）。
func (rg *Registry) Modules() []Module {
	return append([]Module(nil), rg.modules...)
}
