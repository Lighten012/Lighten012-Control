// Package runner 提供模块层执行外部命令的统一入口。
//
// 各模块的命令与参数均为服务端约定的白名单（硬编码在模块代码内，
// 不接受任何外部输入参与拼装）；runner 只负责把参数数组直传 exec
// （不经 shell），并统一超时控制与 stdout / stderr / 退出码收集。
package runner

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// DefaultTimeout 默认命令执行超时。
const DefaultTimeout = 10 * time.Second

// Result 一次命令执行的完整结果。
type Result struct {
	Stdout   string
	Stderr   string
	ExitCode int
}

// Runner 外部命令执行器，可被多个模块并发使用。
type Runner struct {
	timeout time.Duration
}

// New 创建执行器；timeout <= 0 时使用 DefaultTimeout。
func New(timeout time.Duration) *Runner {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}
	return &Runner{timeout: timeout}
}

// Run 执行命令并返回结果。
// 命令无法启动、超时或非零退出码均返回 error（非零退出时错误信息附带 stderr）；
// 此时 Result 仍然返回，供调用方检查已产生的输出。
func (r *Runner) Run(ctx context.Context, name string, args ...string) (Result, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	var stdout, stderr bytes.Buffer
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	// 防止命令的子进程持有管道导致 Wait 长时间阻塞
	cmd.WaitDelay = 2 * time.Second

	err := cmd.Run()
	res := Result{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
	}
	if cmd.ProcessState != nil {
		res.ExitCode = cmd.ProcessState.ExitCode()
	}

	switch {
	case err != nil:
		return res, fmt.Errorf("执行命令 %s 失败: %w", name, err)
	case res.ExitCode != 0:
		return res, fmt.Errorf("命令 %s 退出码 %d: %s", name, res.ExitCode, strings.TrimSpace(res.Stderr))
	}
	return res, nil
}
