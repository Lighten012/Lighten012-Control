# Lighten012-Control

Linux 控制后端：通过 HTTP API 控制 Linux 主机的不同功能模块。
架构自设计之初即按模块分层，容器（Podman）为第一个模块，后续将以相同模式接入
iptables 等其他模块。

## 架构分层

```
HTTP API 层   internal/httpapi   路由组装、统一 JSON 响应、静态页托管、访问日志
模块层        internal/modules   Module 接口 + 注册表；每个 Linux 功能模块一个子包
基础设施层    internal/runner    通用命令执行器（不经 shell、超时控制），所有模块复用
```

前端页面（`web/`，通过 go:embed 编译进二进制）只调用 `/lighten012-api/*` 接口。

## 目录结构

```
cmd/server/main.go            程序入口：配置 → 执行器 → 注册模块 → 启动 HTTP
internal/config/config.go     配置加载（环境变量 + 命令行参数）
internal/runner/runner.go     通用命令执行器（exec 直调，不经 shell）
internal/modules/module.go    Module 接口 + 注册表
internal/modules/podman/      Podman 组件模块（容器/镜像/网络/存储卷/信息 只读查询）
internal/httpapi/             统一响应、路由组装、访问日志
web/                          极简实验页面（后续整体替换）
```

## 环境要求

- Go 1.22+（本机已安装于 `/data/workspace/.tools/go`，可自行调整 PATH）
- podman CLI：服务进程以哪个用户运行，即以该用户的 podman 视图执行命令

## 启动

```bash
make run          # 等价于 go run ./cmd/server
# 或
make build && ./bin/lighten012-control -addr :8080
```

配置项（命令行参数优先于环境变量）：

| 参数 | 环境变量 | 默认值 | 说明 |
| --- | --- | --- | --- |
| `-addr` | `LIGHTEN_ADDR` | `:8080` | HTTP 监听地址 |
| `-cmd-timeout-ms` | `LIGHTEN_CMD_TIMEOUT_MS` | `10000` | 模块命令执行超时（毫秒） |

## API

所有端点均为只读查询，统一响应格式 `{code, message, data}`；路径约定为
`/lighten012-api/<操控组件>/<资源>`。

| 端点 | 内部命令 | 说明 |
| --- | --- | --- |
| GET /lighten012-api/podman/containers | `podman ps -a --format json` | 全部容器（含已停止） |
| GET /lighten012-api/podman/images | `podman images --format json` | 本地镜像 |
| GET /lighten012-api/podman/networks | `podman network ls --format json` | 网络 |
| GET /lighten012-api/podman/volumes | `podman volume ls --format json` | 存储卷 |
| GET /lighten012-api/podman/info | `podman info --format json` | 运行环境信息（版本/存储驱动/CPU/内存/发行版/rootless） |

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
{ "code": 500, "message": "执行 podman ps -a 失败: ...", "data": null }
```

## 新增模块指引（如 iptables）

1. 新建 `internal/modules/iptables/`，实现业务 `Service`：通过
   `runner.Run(ctx, "iptables", args...)` 执行命令（命令与参数在服务端
   硬编码），解析输出并定义 DTO；
2. 实现 `modules.Module` 接口：`Name()` 返回模块名；`RegisterRoutes()` 注册
   自己的 `GET/POST /lighten012-api/...` 路由（根路径用共享常量 `modules.PathPrefix`）；
3. 在 `cmd/server/main.go` 中 `registry.Register(iptables.NewModule(executor))`。

核心路由器、执行器、统一响应格式均无需改动。

## 安全说明（家庭可信内网）

- 仅限家庭可信内网运行，无鉴权，勿暴露公网；路由已按 `/lighten012-api` 统一规划，将来需要时可加认证中间件；
- 各模块的命令均为服务端约定的白名单命令（硬编码在模块代码里，如容器模块固定执行
  `podman ps -a --format json`），参数经 exec 数组直传、不经 shell，不接受任何外部输入参与命令拼装；
- 页面渲染全部使用 `textContent`，不存在 XSS 注入面。
