// Package storage will manage storage.
package storage

import (
	"path/filepath"

	"github.com/dwrth/knaller/internal/config"
)

type Manager struct {
	cfg *config.Config
}

func New(cfg *config.Config) *Manager {
	return &Manager{cfg: cfg}
}

func (m *Manager) BaseRootfs() string {
	return m.cfg.Storage.BaseRootfs
}

func (m *Manager) KernelSource() string {
	return m.cfg.Storage.Kernel
}

func (m *Manager) VMDir(id string) string {
	return filepath.Join(m.cfg.Firecracker.VmDirectory, id)
}

func (m *Manager) Rootfs(id string) string {
	return filepath.Join(m.VMDir(id), "rootfs.ext4")
}

func (m *Manager) Kernel(id string) string {
	return filepath.Join(m.VMDir(id), "vmlinux")
}

func (m *Manager) Config(id string) string {
	return filepath.Join(m.VMDir(id), "config.json")
}
