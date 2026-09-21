package podman

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

// unmarshal 把样例 JSON 反序列化为 T，模拟 query 的解析环节，
// 从而同时覆盖 json tag 与映射函数两层正确性。
func unmarshal[T any](t *testing.T, data string) T {
	t.Helper()
	var out T
	if err := json.Unmarshal([]byte(data), &out); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	return out
}

func TestMapContainer(t *testing.T) {
	raw := unmarshal[[]podmanContainer](t, `[
		{
			"Id": "e90d0b15ed0712345678901234567890abcdef1234567890abcdef1234567890",
			"Names": ["web"],
			"Image": "docker.io/library/alpine:latest",
			"State": "running",
			"Status": "Up 5 minutes",
			"CreatedAt": "5 minutes ago",
			"Created": 1789677440
		},
		{
			"Id": "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
			"Names": ["db"],
			"State": "created",
			"Status": "Created",
			"Created": 1789677440
		}
	]`)

	got := mapAll(raw, mapContainer)
	want := []Container{
		{ID: "e90d0b15ed07", Name: "web", Image: "docker.io/library/alpine:latest", State: "running", Status: "Up 5 minutes", Created: "5 minutes ago"},
		{ID: "0123456789ab", Name: "db", State: "created", Status: "Created", Created: time.Unix(1789677440, 0).Format("2006-01-02 15:04:05")},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("mapAll(mapContainer) =\n%+v\nwant\n%+v", got, want)
	}
}

// 样例取自 podman 5.8.4 实测输出（RepoTags 为 null、名称取 Names[0]）。
func TestMapImage(t *testing.T) {
	raw := unmarshal[[]podmanImage](t, `[
		{
			"Id": "320994c3b997e2ec6433f717f153e108023c5bec8fefa8d76b83451d16d05ea8",
			"ParentId": "",
			"RepoTags": null,
			"RepoDigests": ["docker.io/library/alpine@sha256:294b683cb724975bec92580e1e685676bd4b50bda910ddb8c51d4cabeaec77e6"],
			"Size": 8715873,
			"SharedSize": 0,
			"VirtualSize": 8715873,
			"Labels": null,
			"Containers": 2,
			"Digest": "sha256:294b683cb724975bec92580e1e685676bd4b50bda910ddb8c51d4cabeaec77e6",
			"Names": ["docker.io/library/alpine:latest"],
			"Created": 1789677440,
			"CreatedAt": "2026-09-17T20:37:20Z"
		}
	]`)

	got := mapAll(raw, mapImage)
	want := []Image{
		{
			ID:         "320994c3b997",
			Name:       "docker.io/library/alpine:latest",
			Size:       8715873,
			Containers: 2,
			Created:    "2026-09-17T20:37:20Z",
			Digest:     "sha256:294b683cb724975bec92580e1e685676bd4b50bda910ddb8c51d4cabeaec77e6",
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("mapAll(mapImage) =\n%+v\nwant\n%+v", got, want)
	}
}

// 样例取自 podman 5.8.4 实测输出（键全小写）。
func TestMapNetwork(t *testing.T) {
	raw := unmarshal[[]podmanNetwork](t, `[
		{
			"name": "podman",
			"id": "2f259bab93aaaaa2542ba43ef33eb990d0999ee1b9924b557b7be53c0b7a1bb9",
			"driver": "bridge",
			"network_interface": "podman0",
			"created": "2026-09-20T17:04:52.016650082+08:00",
			"subnets": [{"subnet": "10.88.0.0/16", "gateway": "10.88.0.1"}],
			"ipv6_enabled": false,
			"internal": false,
			"dns_enabled": false
		}
	]`)

	got := mapAll(raw, mapNetwork)
	want := []Network{
		{
			ID:          "2f259bab93aa",
			Name:        "podman",
			Driver:      "bridge",
			Interface:   "podman0",
			Created:     "2026-09-20T17:04:52.016650082+08:00",
			Subnets:     []string{"10.88.0.0/16"},
			DNSEnabled:  false,
			Internal:    false,
			IPV6Enabled: false,
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("mapAll(mapNetwork) =\n%+v\nwant\n%+v", got, want)
	}
}

// 卷字段当前环境无法实测（无卷），样例按 DTO 设计的键名编写，
// 验证阶段会以临时卷实测校准。
func TestMapVolume(t *testing.T) {
	raw := unmarshal[[]podmanVolume](t, `[
		{
			"Name": "lighten-test-vol",
			"Driver": "local",
			"Scope": "local",
			"Mountpoint": "/data/containers/storage/volumes/lighten-test-vol/_data",
			"CreatedAt": "2026-09-20T17:30:00+08:00",
			"Anonymous": false
		}
	]`)

	got := mapAll(raw, mapVolume)
	want := []Volume{
		{
			Name:       "lighten-test-vol",
			Driver:     "local",
			Scope:      "local",
			Mountpoint: "/data/containers/storage/volumes/lighten-test-vol/_data",
			Created:    "2026-09-20T17:30:00+08:00",
			Anonymous:  false,
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("mapAll(mapVolume) =\n%+v\nwant\n%+v", got, want)
	}
}

// 样例取自 podman 5.8.4 实测输出（host 段全小写键、store.graphDriverName、
// version.Version 大写 V）。
func TestMapInfo(t *testing.T) {
	raw := unmarshal[podmanInfo](t, `{
		"host": {
			"arch": "amd64",
			"os": "linux",
			"kernel": "6.6.119-49.27.tl4.x86_64",
			"cpus": 8,
			"memTotal": 16231919616,
			"distribution": {"distribution": "tencentos", "version": "4.4"},
			"security": {"rootless": false}
		},
		"store": {"graphDriverName": "vfs"},
		"version": {"Version": "5.8.4"}
	}`)

	got := mapInfo(raw)
	want := Info{
		Version:       "5.8.4",
		StorageDriver: "vfs",
		Os:            "linux",
		Arch:          "amd64",
		Kernel:        "6.6.119-49.27.tl4.x86_64",
		Cpus:          8,
		MemTotal:      16231919616,
		Rootless:      false,
		Distro:        "tencentos",
		DistroVersion: "4.4",
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("mapInfo() =\n%+v\nwant\n%+v", got, want)
	}
}
