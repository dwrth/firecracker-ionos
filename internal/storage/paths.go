// Package storage provisions per-sandbox disks and kernels
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

// SandboxDir returns the directory for the given sandbox ID.
func (m *Manager) SandboxDir(id string) string {
	return filepath.Join(m.cfg.Firecracker.SandboxDirectory, id)
}

// Rootfs returns the rootfs path for the given sandbox ID.
func (m *Manager) Rootfs(id string) string {
	return filepath.Join(m.SandboxDir(id), "rootfs.ext4")
}

// Kernel returns the kernel path for the given sandbox ID.
func (m *Manager) Kernel(id string) string {
	return filepath.Join(m.SandboxDir(id), "vmlinux")
}

// Config returns the Firecracker config path for the given sandbox ID.
func (m *Manager) Config(id string) string {
	return filepath.Join(m.SandboxDir(id), "config.json")
}
