// Package storage provisions per-VM disks and kernels
package storage

import (
	"path/filepath"

	"github.com/dwrth/knaller/internal/config"
)

// Manager resolves storage paths from configuration.
type Manager struct {
	cfg *config.Config
}

// New returns a Manager for cfg.
func New(cfg *config.Config) *Manager {
	return &Manager{cfg: cfg}
}

// BaseRootfs returns the path to the shared base rootfs image.
func (m *Manager) BaseRootfs() string {
	return m.cfg.Storage.BaseRootfs
}

// KernelSource returns the path to the shared kernel image.
func (m *Manager) KernelSource() string {
	return m.cfg.Storage.Kernel
}

// VMDir returns the directory for VM id.
func (m *Manager) VMDir(id string) string {
	return filepath.Join(m.cfg.Firecracker.VmDirectory, id)
}

// Rootfs returns the rootfs path for VM id.
func (m *Manager) Rootfs(id string) string {
	return filepath.Join(m.VMDir(id), "rootfs.ext4")
}

// Kernel returns the kernel path for VM id.
func (m *Manager) Kernel(id string) string {
	return filepath.Join(m.VMDir(id), "vmlinux")
}

// Config returns the Firecracker config path for VM id.
func (m *Manager) Config(id string) string {
	return filepath.Join(m.VMDir(id), "config.json")
}
