package podman

import (
	"net"
	"reflect"
	"testing"
	"time"

	"github.com/containers/podman/v5/libpod/define"
	bindingsTypes "github.com/containers/podman/v5/pkg/domain/entities/types"
	netTypes "go.podman.io/common/libnetwork/types"
)

func TestMapContainer(t *testing.T) {
	created := time.Unix(1789677440, 0)
	raw := []bindingsTypes.ListContainer{
		{
			ID:        "e90d0b15ed0712345678901234567890abcdef1234567890abcdef1234567890",
			Names:     []string{"web"},
			Image:     "docker.io/library/alpine:latest",
			State:     "running",
			Status:    "Up 5 minutes",
			Created:   created,
			CreatedAt: "5 minutes ago",
		},
		{
			ID:      "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
			Names:   []string{"db"},
			State:   "created",
			Status:  "Created",
			Created: created,
		},
	}

	got := mapAll(raw, mapContainer)
	want := []Container{
		{ID: "e90d0b15ed07", Name: "web", Image: "docker.io/library/alpine:latest", State: "running", Status: "Up 5 minutes", Created: "5 minutes ago"},
		{ID: "0123456789ab", Name: "db", State: "created", Status: "Created", Created: created.Local().Format("2006-01-02 15:04:05")},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("mapAll(mapContainer) =\n%+v\nwant\n%+v", got, want)
	}
}

func TestMapImage(t *testing.T) {
	created := int64(1789677440)
	raw := []*bindingsTypes.ImageSummary{
		{
			ID:          "320994c3b997e2ec6433f717f153e108023c5bec8fefa8d76b83451d16d05ea8",
			Names:       []string{"docker.io/library/alpine:latest"},
			RepoDigests: []string{"docker.io/library/alpine@sha256:294b683cb724975bec92580e1e685676bd4b50bda910ddb8c51d4cabeaec77e6"},
			Digest:      "sha256:294b683cb724975bec92580e1e685676bd4b50bda910ddb8c51d4cabeaec77e6",
			Size:        8715873,
			Containers:  2,
			Created:     created,
		},
	}

	got := make([]Image, 0, len(raw))
	for _, item := range raw {
		got = append(got, mapImage(item))
	}
	want := []Image{
		{
			ID:         "320994c3b997",
			Name:       "docker.io/library/alpine:latest",
			Size:       8715873,
			Containers: 2,
			Created:    time.Unix(created, 0).UTC().Format(time.RFC3339),
			Digest:     "sha256:294b683cb724975bec92580e1e685676bd4b50bda910ddb8c51d4cabeaec77e6",
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("mapImage() =\n%+v\nwant\n%+v", got, want)
	}
}

func TestMapNetwork(t *testing.T) {
	subnet, err := netTypes.ParseCIDR("10.88.0.0/16")
	if err != nil {
		t.Fatalf("netTypes.ParseCIDR() error = %v", err)
	}
	created := time.Unix(1789677440, 0)
	raw := []netTypes.Network{
		{
			Name:             "podman",
			ID:               "2f259bab93aaaaa2542ba43ef33eb990d0999ee1b9924b557b7be53c0b7a1bb9",
			Driver:           "bridge",
			NetworkInterface: "podman0",
			Created:          created,
			Subnets:          []netTypes.Subnet{{Subnet: subnet, Gateway: net.ParseIP("10.88.0.1")}},
			IPv6Enabled:      false,
			Internal:         false,
			DNSEnabled:       false,
		},
	}

	got := mapAll(raw, mapNetwork)
	want := []Network{
		{
			ID:          "2f259bab93aa",
			Name:        "podman",
			Driver:      "bridge",
			Interface:   "podman0",
			Created:     created.Local().Format("2006-01-02 15:04:05"),
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

func TestMapVolume(t *testing.T) {
	created := time.Unix(1789677440, 0)
	raw := []*bindingsTypes.VolumeListReport{
		{
			VolumeConfigResponse: bindingsTypes.VolumeConfigResponse{
				InspectVolumeData: define.InspectVolumeData{
					Name:       "lighten-test-vol",
					Driver:     "local",
					Scope:      "local",
					Mountpoint: "/data/containers/storage/volumes/lighten-test-vol/_data",
					CreatedAt:  created,
					Anonymous:  false,
				},
			},
		},
	}

	got := make([]Volume, 0, len(raw))
	for _, item := range raw {
		got = append(got, mapVolume(item))
	}
	want := []Volume{
		{
			Name:       "lighten-test-vol",
			Driver:     "local",
			Scope:      "local",
			Mountpoint: "/data/containers/storage/volumes/lighten-test-vol/_data",
			Created:    created.Local().Format("2006-01-02 15:04:05"),
			Anonymous:  false,
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("mapVolume() =\n%+v\nwant\n%+v", got, want)
	}
}

func TestMapInfo(t *testing.T) {
	raw := &define.Info{
		Host: &define.HostInfo{
			Arch:     "amd64",
			OS:       "linux",
			Kernel:   "6.6.119-49.27.tl4.x86_64",
			CPUs:     8,
			MemTotal: 16231919616,
			Distribution: define.DistributionInfo{
				Distribution: "tencentos",
				Version:      "4.4",
			},
			Security: define.SecurityInfo{
				Rootless: false,
			},
		},
		Store: &define.StoreInfo{
			GraphDriverName: "vfs",
		},
		Version: define.Version{
			Version: "5.8.4",
		},
	}

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
