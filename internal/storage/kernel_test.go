package storage_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/dwrth/knaller/internal/storage"
)

func TestInstallKernelMissingSource(t *testing.T) {
	tmp := t.TempDir()
	cfg := testConfig()
	cfg.Storage.Kernel = filepath.Join(tmp, "missing-vmlinux")
	cfg.Firecracker.VmDirectory = filepath.Join(tmp, "vms")

	m := storage.New(cfg)
	err := m.InstallKernel("vm-0001")
	if err == nil {
		t.Fatal("expected error for missing kernel source")
	}
}

func TestInstallKernelAlreadyExists(t *testing.T) {
	tmp := t.TempDir()
	cfg := testConfig()
	cfg.Storage.Kernel = filepath.Join(tmp, "vmlinux")
	cfg.Firecracker.VmDirectory = filepath.Join(tmp, "vms")

	if err := os.WriteFile(cfg.Storage.Kernel, []byte("kernel"), 0o644); err != nil {
		t.Fatal(err)
	}

	m := storage.New(cfg)
	kernel := m.Kernel("vm-0001")
	if err := os.MkdirAll(filepath.Dir(kernel), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(kernel, []byte("existing"), 0o644); err != nil {
		t.Fatal(err)
	}

	err := m.InstallKernel("vm-0001")
	if err == nil {
		t.Fatal("expected error when kernel already exists")
	}
}

func TestInstallKernel(t *testing.T) {
	tmp := t.TempDir()
	cfg := testConfig()
	cfg.Storage.Kernel = filepath.Join(tmp, "vmlinux")
	cfg.Firecracker.VmDirectory = filepath.Join(tmp, "vms")

	want := []byte("pinned-kernel-bytes")
	if err := os.WriteFile(cfg.Storage.Kernel, want, 0o644); err != nil {
		t.Fatal(err)
	}

	m := storage.New(cfg)
	if err := m.InstallKernel("vm-0001"); err != nil {
		t.Fatal(err)
	}

	got, err := os.ReadFile(m.Kernel("vm-0001"))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(want) {
		t.Fatalf("installed kernel = %q, want %q", got, want)
	}
}
