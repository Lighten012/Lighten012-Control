package podman

import (
	"strings"
	"time"

	"github.com/containers/podman/v5/libpod/define"
	bindingsTypes "github.com/containers/podman/v5/pkg/domain/entities/types"
	netTypes "go.podman.io/common/libnetwork/types"
)

// ---------- DTO：对外暴露的稳定结构 ----------

// Container 容器信息。
type Container struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Image   string `json:"image"`
	State   string `json:"state"`
	Status  string `json:"status"`
	Created string `json:"created"`
}

// Image 镜像信息。
type Image struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Size       int64  `json:"size"`
	Containers int    `json:"containers"`
	Created    string `json:"created"`
	Digest     string `json:"digest"`
}

// Network 网络信息。
type Network struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Driver      string   `json:"driver"`
	Interface   string   `json:"interface"`
	Created     string   `json:"created"`
	Subnets     []string `json:"subnets"`
	DNSEnabled  bool     `json:"dnsEnabled"`
	Internal    bool     `json:"internal"`
	IPV6Enabled bool     `json:"ipv6Enabled"`
}

// Volume 存储卷信息。
type Volume struct {
	Name       string `json:"name"`
	Driver     string `json:"driver"`
	Scope      string `json:"scope"`
	Mountpoint string `json:"mountpoint"`
	Created    string `json:"created"`
	Anonymous  bool   `json:"anonymous"`
}

// Info 运行环境信息。
type Info struct {
	Version       string `json:"version"`
	StorageDriver string `json:"storageDriver"`
	Os            string `json:"os"`
	Arch          string `json:"arch"`
	Kernel        string `json:"kernel"`
	Cpus          int    `json:"cpus"`
	MemTotal      int64  `json:"memTotal"`
	Rootless      bool   `json:"rootless"`
	Distro        string `json:"distro"`
	DistroVersion string `json:"distroVersion"`
}

// ---------- 映射 ----------

// shortID 截断 ID 为 12 位短 ID（与 podman CLI 默认展示一致）。
func shortID(id string) string {
	id = strings.TrimSpace(id)
	if len(id) > 12 {
		return id[:12]
	}
	return id
}

// displayTime 优先返回 API 提供的人性化时间，否则格式化时间戳。
func displayTime(created time.Time, createdAt string) string {
	if createdAt != "" {
		return createdAt
	}
	if !created.IsZero() {
		return created.Local().Format("2006-01-02 15:04:05")
	}
	return ""
}

// formatUnix 将 Podman 返回的 Unix 时间戳格式化为稳定展示格式。
func formatUnix(seconds int64) string {
	if seconds == 0 {
		return ""
	}
	return time.Unix(seconds, 0).UTC().Format(time.RFC3339)
}

// firstName 依次取候选数组的首个非空元素，全部为空时返回 "<none>"。
func firstName(candidates ...[]string) string {
	for _, list := range candidates {
		for _, s := range list {
			if strings.TrimSpace(s) != "" {
				return s
			}
		}
	}
	return "<none>"
}

func mapContainer(raw bindingsTypes.ListContainer) Container {
	c := Container{
		ID:      shortID(raw.ID),
		Image:   raw.Image,
		State:   raw.State,
		Status:  raw.Status,
		Created: displayTime(raw.Created, raw.CreatedAt),
	}
	if len(raw.Names) > 0 {
		c.Name = raw.Names[0]
	}
	return c
}

func mapImage(raw *bindingsTypes.ImageSummary) Image {
	return Image{
		ID:         shortID(raw.ID),
		Name:       firstName(raw.Names, raw.RepoTags),
		Size:       raw.Size,
		Containers: raw.Containers,
		Created:    formatUnix(raw.Created),
		Digest:     raw.Digest,
	}
}

func mapNetwork(raw netTypes.Network) Network {
	subnets := make([]string, 0, len(raw.Subnets))
	for _, s := range raw.Subnets {
		subnets = append(subnets, s.Subnet.String())
	}
	return Network{
		ID:          shortID(raw.ID),
		Name:        raw.Name,
		Driver:      raw.Driver,
		Interface:   raw.NetworkInterface,
		Created:     displayTime(raw.Created, ""),
		Subnets:     subnets,
		DNSEnabled:  raw.DNSEnabled,
		Internal:    raw.Internal,
		IPV6Enabled: raw.IPv6Enabled,
	}
}

func mapVolume(raw *bindingsTypes.VolumeListReport) Volume {
	return Volume{
		Name:       raw.Name,
		Driver:     raw.Driver,
		Scope:      raw.Scope,
		Mountpoint: raw.Mountpoint,
		Created:    displayTime(raw.CreatedAt, ""),
		Anonymous:  raw.Anonymous,
	}
}

func mapInfo(raw *define.Info) Info {
	if raw == nil {
		return Info{}
	}

	info := Info{}
	if raw.Host != nil {
		info.Os = raw.Host.OS
		info.Arch = raw.Host.Arch
		info.Kernel = raw.Host.Kernel
		info.Cpus = raw.Host.CPUs
		info.MemTotal = raw.Host.MemTotal
		info.Rootless = raw.Host.Security.Rootless
		info.Distro = raw.Host.Distribution.Distribution
		info.DistroVersion = raw.Host.Distribution.Version
	}
	if raw.Store != nil {
		info.StorageDriver = raw.Store.GraphDriverName
	}
	info.Version = raw.Version.Version
	return info
}
