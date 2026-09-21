package podman

import (
	"strings"
	"time"
)

// ---------- 原始结构：podman --format json 的实测字段（未列出的自动忽略） ----------

// podmanContainer podman ps -a --format json 条目。
// podman 5.x 的 Created 为 Unix 时间戳数值（旧版可能为字符串），
// CreatedAt 为人性化时间字符串，故 Created 用 any 兼容两种形态。
type podmanContainer struct {
	Id        string   `json:"Id"`
	Names     []string `json:"Names"`
	Image     string   `json:"Image"`
	State     string   `json:"State"`
	Status    string   `json:"Status"`
	Created   any      `json:"Created"`
	CreatedAt string   `json:"CreatedAt"`
}

// podmanImage podman images --format json 条目（实测：Id 为裸 ID，
// RepoTags 可能为 null，名称优先取 Names[0]）。
type podmanImage struct {
	Id          string   `json:"Id"`
	Names       []string `json:"Names"`
	RepoTags    []string `json:"RepoTags"`
	RepoDigests []string `json:"RepoDigests"`
	Digest      string   `json:"Digest"`
	Size        int64    `json:"Size"`
	Containers  int      `json:"Containers"`
	CreatedAt   string   `json:"CreatedAt"`
}

// podmanNetwork podman network ls --format json 条目（实测键全小写）。
type podmanNetwork struct {
	Name             string `json:"name"`
	Id               string `json:"id"`
	Driver           string `json:"driver"`
	NetworkInterface string `json:"network_interface"`
	Created          string `json:"created"`
	Subnets          []struct {
		Subnet  string `json:"subnet"`
		Gateway string `json:"gateway"`
	} `json:"subnets"`
	IPV6Enabled bool `json:"ipv6_enabled"`
	Internal    bool `json:"internal"`
	DNSEnabled  bool `json:"dns_enabled"`
}

// podmanVolume podman volume ls --format json 条目。
// 当前环境暂无卷，键名按 podman 惯例设计，验证阶段以实测校准。
type podmanVolume struct {
	Name       string `json:"Name"`
	Driver     string `json:"Driver"`
	Scope      string `json:"Scope"`
	Mountpoint string `json:"Mountpoint"`
	CreatedAt  string `json:"CreatedAt"`
	Anonymous  bool   `json:"Anonymous"`
}

// podmanInfo podman info --format json 中所需子集（实测：host 段全小写键，
// store.graphDriverName，version.Version）。
type podmanInfo struct {
	Host struct {
		Arch         string `json:"arch"`
		Os           string `json:"os"`
		Kernel       string `json:"kernel"`
		Cpus         int    `json:"cpus"`
		MemTotal     int64  `json:"memTotal"`
		Distribution struct {
			Distribution string `json:"distribution"`
			Version      string `json:"version"`
		} `json:"distribution"`
		Security struct {
			Rootless bool `json:"rootless"`
		} `json:"security"`
	} `json:"host"`
	Store struct {
		GraphDriverName string `json:"graphDriverName"`
	} `json:"store"`
	Version struct {
		Version string `json:"Version"`
	} `json:"version"`
}

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

// displayTime 计算容器创建时间展示值：优先 CreatedAt（人性化字符串），
// 缺失时回退格式化 Created（Unix 时间戳或字符串）。
func displayTime(created any, createdAt string) string {
	if createdAt != "" {
		return createdAt
	}
	switch v := created.(type) {
	case float64:
		if v > 0 {
			return time.Unix(int64(v), 0).Format("2006-01-02 15:04:05")
		}
	case string:
		return v
	}
	return ""
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

func mapContainer(raw podmanContainer) Container {
	c := Container{
		ID:      shortID(raw.Id),
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

func mapImage(raw podmanImage) Image {
	return Image{
		ID:         shortID(raw.Id),
		Name:       firstName(raw.Names, raw.RepoTags),
		Size:       raw.Size,
		Containers: raw.Containers,
		Created:    raw.CreatedAt,
		Digest:     raw.Digest,
	}
}

func mapNetwork(raw podmanNetwork) Network {
	subnets := make([]string, 0, len(raw.Subnets))
	for _, s := range raw.Subnets {
		subnets = append(subnets, s.Subnet)
	}
	return Network{
		ID:          shortID(raw.Id),
		Name:        raw.Name,
		Driver:      raw.Driver,
		Interface:   raw.NetworkInterface,
		Created:     raw.Created,
		Subnets:     subnets,
		DNSEnabled:  raw.DNSEnabled,
		Internal:    raw.Internal,
		IPV6Enabled: raw.IPV6Enabled,
	}
}

func mapVolume(raw podmanVolume) Volume {
	return Volume{
		Name:       raw.Name,
		Driver:     raw.Driver,
		Scope:      raw.Scope,
		Mountpoint: raw.Mountpoint,
		Created:    raw.CreatedAt,
		Anonymous:  raw.Anonymous,
	}
}

func mapInfo(raw podmanInfo) Info {
	return Info{
		Version:       raw.Version.Version,
		StorageDriver: raw.Store.GraphDriverName,
		Os:            raw.Host.Os,
		Arch:          raw.Host.Arch,
		Kernel:        raw.Host.Kernel,
		Cpus:          raw.Host.Cpus,
		MemTotal:      raw.Host.MemTotal,
		Rootless:      raw.Host.Security.Rootless,
		Distro:        raw.Host.Distribution.Distribution,
		DistroVersion: raw.Host.Distribution.Version,
	}
}
