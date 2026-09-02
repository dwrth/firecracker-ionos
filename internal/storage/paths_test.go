package storage_test

import (
	"testing"

	"github.com/dwrth/knaller/internal/config"
	"github.com/dwrth/knaller/internal/storage"
)

func testConfig() *config.Config {
	return &config.Config{
		Storage: config.Storage{
			BaseRootfs: "/var/lib/knaller/images/base-rootfs.ext4",
			Kernel:     "/var/lib/knaller/images/vmlinux-6.18.44",
		},
		Firecracker: config.Firecracker{
			SandboxDirectory: "/var/lib/knaller/vms",
		},
	}
}

func TestManagerPaths(t *testing.T) {
	m := storage.New(testConfig())

	if got, want := m.BaseRootfs(), "/var/lib/knaller/images/base-rootfs.ext4"; got != want {
		t.Errorf("BaseRootfs() = %q, want %q", got, want)
	}
	if got, want := m.KernelSource(), "/var/lib/knaller/images/vmlinux-6.18.44"; got != want {
		t.Errorf("KernelSource() = %q, want %q", got, want)
	}
	if got, want := m.SandboxDir("vm-0004"), "/var/lib/knaller/vms/vm-0004"; got != want {
		t.Errorf("SandboxDir() = %q, want %q", got, want)
	}
	if got, want := m.Rootfs("vm-0004"), "/var/lib/knaller/vms/vm-0004/rootfs.ext4"; got != want {
		t.Errorf("Rootfs() = %q, want %q", got, want)
	}
	if got, want := m.Kernel("vm-0004"), "/var/lib/knaller/vms/vm-0004/vmlinux"; got != want {
		t.Errorf("Kernel() = %q, want %q", got, want)
	}
	if got, want := m.Config("vm-0004"), "/var/lib/knaller/vms/vm-0004/config.json"; got != want {
		t.Errorf("Config() = %q, want %q", got, want)
	}
}
