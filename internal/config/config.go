// Package config 提供服务运行配置的加载。
package config

import (
	"flag"
	"log"
	"os"
	"strconv"
	"time"
)

const (
	EnvAddr         = "LIGHTEN_ADDR"
	EnvCmdTimeoutMs = "LIGHTEN_CMD_TIMEOUT_MS"

	DefaultAddr       = ":8080"
	DefaultCmdTimeout = 10 * time.Second
)

// Settings 服务运行配置。
type Settings struct {
	Addr       string        // HTTP 监听地址，例如 ":8080"
	CmdTimeout time.Duration // 模块命令执行超时
}

// Load 加载配置：环境变量提供默认值，命令行参数可覆盖。
func Load() Settings {
	envAddr := os.Getenv(EnvAddr)
	if envAddr == "" {
		envAddr = DefaultAddr
	}

	envTimeout := DefaultCmdTimeout
	if raw := os.Getenv(EnvCmdTimeoutMs); raw != "" {
		if n, err := strconv.Atoi(raw); err == nil && n > 0 {
			envTimeout = time.Duration(n) * time.Millisecond
		} else {
			log.Printf("[config] 忽略非法的 %s=%q", EnvCmdTimeoutMs, raw)
		}
	}

	addr := flag.String("addr", envAddr, "HTTP 监听地址")
	timeoutMs := flag.Int("cmd-timeout-ms", int(envTimeout.Milliseconds()), "模块命令执行超时（毫秒）")
	flag.Parse()

	return Settings{
		Addr:       *addr,
		CmdTimeout: time.Duration(*timeoutMs) * time.Millisecond,
	}
}
