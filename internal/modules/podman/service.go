// Package podman Podman 组件模块：以只读方式查询容器、镜像、网络、
// 存储卷与运行环境信息。
//
// 模块通过 Podman 官方 Go bindings 访问本机 Podman system service，
// 并把 bindings 返回的数据映射为 DTO；原始结构不会泄漏到模块层之外。
package podman

import (
	"context"
	"fmt"
	"time"

	"github.com/containers/podman/v5/pkg/bindings"
	bindingsContainers "github.com/containers/podman/v5/pkg/bindings/containers"
	bindingsImages "github.com/containers/podman/v5/pkg/bindings/images"
	bindingsNetwork "github.com/containers/podman/v5/pkg/bindings/network"
	bindingsSystem "github.com/containers/podman/v5/pkg/bindings/system"
	bindingsVolumes "github.com/containers/podman/v5/pkg/bindings/volumes"
)

// Service Podman 组件的业务逻辑。
type Service struct {
	uri     string
	timeout time.Duration
}

// NewService 创建 Podman 组件服务。
func NewService(uri string, timeout time.Duration) *Service {
	if timeout <= 0 {
		timeout = 10 * time.Second
	}
	return &Service{uri: uri, timeout: timeout}
}

func (s *Service) connection(ctx context.Context) (context.Context, context.CancelFunc, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	conn, err := bindings.NewConnection(ctx, s.uri)
	if err != nil {
		cancel()
		return nil, nil, fmt.Errorf("连接 Podman socket 失败: %w", err)
	}
	return conn, cancel, nil
}

func mapAll[F, T any](raw []F, mapFn func(F) T) []T {
	list := make([]T, 0, len(raw))
	for _, item := range raw {
		list = append(list, mapFn(item))
	}
	return list
}

// ListContainers 列出全部容器（含已停止）。
func (s *Service) ListContainers(ctx context.Context) ([]Container, error) {
	conn, cancel, err := s.connection(ctx)
	if err != nil {
		return nil, err
	}
	defer cancel()

	raw, err := bindingsContainers.List(conn, new(bindingsContainers.ListOptions).WithAll(true))
	if err != nil {
		return nil, fmt.Errorf("查询容器列表失败: %w", err)
	}
	return mapAll(raw, mapContainer), nil
}

// ListImages 列出本地镜像。
func (s *Service) ListImages(ctx context.Context) ([]Image, error) {
	conn, cancel, err := s.connection(ctx)
	if err != nil {
		return nil, err
	}
	defer cancel()

	raw, err := bindingsImages.List(conn, nil)
	if err != nil {
		return nil, fmt.Errorf("查询镜像列表失败: %w", err)
	}

	list := make([]Image, 0, len(raw))
	for _, item := range raw {
		if item != nil {
			list = append(list, mapImage(item))
		}
	}
	return list, nil
}

// ListNetworks 列出网络。
func (s *Service) ListNetworks(ctx context.Context) ([]Network, error) {
	conn, cancel, err := s.connection(ctx)
	if err != nil {
		return nil, err
	}
	defer cancel()

	raw, err := bindingsNetwork.List(conn, nil)
	if err != nil {
		return nil, fmt.Errorf("查询网络列表失败: %w", err)
	}
	return mapAll(raw, mapNetwork), nil
}

// ListVolumes 列出存储卷。
func (s *Service) ListVolumes(ctx context.Context) ([]Volume, error) {
	conn, cancel, err := s.connection(ctx)
	if err != nil {
		return nil, err
	}
	defer cancel()

	raw, err := bindingsVolumes.List(conn, nil)
	if err != nil {
		return nil, fmt.Errorf("查询存储卷列表失败: %w", err)
	}

	list := make([]Volume, 0, len(raw))
	for _, item := range raw {
		if item != nil {
			list = append(list, mapVolume(item))
		}
	}
	return list, nil
}

// Info 查询运行环境信息。
func (s *Service) Info(ctx context.Context) (Info, error) {
	conn, cancel, err := s.connection(ctx)
	if err != nil {
		return Info{}, err
	}
	defer cancel()

	raw, err := bindingsSystem.Info(conn, nil)
	if err != nil {
		return Info{}, fmt.Errorf("查询运行环境信息失败: %w", err)
	}
	return mapInfo(raw), nil
}
