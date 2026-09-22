# Lighten012-Control

Linux 控制后端：通过 HTTP API 控制 Linux 主机的不同功能模块。
架构按模块分层，容器（Podman）为当前模块，路由组合在服务入口显式完成；
Podman 查询通过官方 Go bindings 访问 Podman system service，不再调用 podman CLI。

## 架构分层

```
HTTP API 层   internal/httpapi   统一 JSON 响应、路由辅助、静态页托管、访问日志
模块层        internal/modules   每个 Linux 功能模块一个子包
```

前端页面（`web/`，通过 go:embed 编译进二进制）只调用 `/lighten012-api/*` 接口。

## 目录结构

```
cmd/server/main.go            程序入口：检查 Podman socket → 组装模块路由 → 启动 HTTP
internal/modules/podman/      Podman 组件模块（容器/镜像/网络/存储卷/信息 只读查询）
internal/httpapi/             统一响应、路由组装、访问日志
web/                          极简实验页面（后续整体替换）
```

## 环境要求

- Go（以 `go.mod` 中的版本要求为准；本机可安装于 `/data/workspace/.tools/go`）
- Podman 5.x 与 `podman.socket`

rootful 服务：

```bash
systemctl enable --now podman.socket
```

rootless 服务：

```bash
systemctl --user enable --now podman.socket
```

服务进程以哪个用户运行，就使用该用户可见的 Podman socket。默认按 rootful
`/run/podman/podman.sock` 或 rootless `$XDG_RUNTIME_DIR/podman/podman.sock`
自动选择；也可用环境变量覆盖：

```bash
PODMAN_URI=unix:///run/user/1000/podman/podman.sock make run
```

## 启动

```bash
make run
# 或
make build && ./bin/lighten012-control
```

## API

所有端点均为只读查询，统一响应格式 `{code, message, data}`；路径约定为
`/lighten012-api/<操控组件>/<资源>`。

| 端点 | 说明 |
| --- | --- |
| GET /lighten012-api/podman/containers | 全部容器（含已停止） |
| GET /lighten012-api/podman/images | 本地镜像 |
| GET /lighten012-api/podman/networks | 网络 |
| GET /lighten012-api/podman/volumes | 存储卷 |
| GET /lighten012-api/podman/info | 运行环境信息（版本/存储驱动/CPU/内存/发行版/rootless） |

容器列表示例响应：

```json
{
  "code": 0,
  "message": "ok",
  "data": [
    {
      "id": "e90d0b15ed07",
      "name": "web",
      "image": "docker.io/library/alpine:latest",
      "state": "running",
      "status": "Up 5 minutes",
      "created": "5 minutes ago"
    }
  ]
}
```

错误响应：

```json
{ "code": 500, "message": "查询容器列表失败: ...", "data": null }
```

## 新增模块指引

1. 新建 `internal/modules/<新模块>/`，实现业务 `Service`：通过该模块可长期
   使用的 Go API 或系统能力完成查询，并定义稳定的 DTO；
2. 提供 `NewModule(...)` 和 `RegisterRoutes(mux *http.ServeMux)`，
   路由使用共享根路径 `httpapi.PathPrefix`；
3. 在 `cmd/server/main.go` 中显式调用：
   `newModule(...).RegisterRoutes(mux)`。

## 安全说明（家庭可信内网）

- 仅限家庭可信内网运行，无鉴权，勿暴露公网；路由已按 `/lighten012-api` 统一规划，将来需要时可加认证中间件；
- 当前仅暴露只读查询；Podman 模块只通过本机 Unix socket 访问 Podman system service，不拼接命令，不接受外部输入参与查询参数拼装；
- 页面渲染全部使用 `textContent`，不存在 XSS 注入面。

## 构建标签

`make build/test` 会使用 `remote exclude_graphdriver_btrfs containers_image_openpgp`。
本项目只作为 Podman REST bindings 客户端，这些标签可避免引入不需要的本地存储
C 依赖。直接执行 `go test ./...` 时需要附加相同标签。
