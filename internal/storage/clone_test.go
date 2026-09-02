package storage_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/dwrth/knaller/internal/storage"
)

func TestCloneBaseRootfsMissingBase(t *testing.T) {
	tmp := t.TempDir()
	cfg := testConfig()
	cfg.Storage.BaseRootfs = filepath.Join(tmp, "missing.ext4")
	cfg.Firecracker.SandboxDirectory = filepath.Join(tmp, "vms")

	m := storage.New(cfg)
	err := m.CloneBaseRootfs("vm-0001")
	if err == nil {
		t.Fatal("expected error for missing base rootfs")
	}
}

func TestCloneBaseRootfsAlreadyExists(t *testing.T) {
	tmp := t.TempDir()
	cfg := testConfig()
	cfg.Storage.BaseRootfs = filepath.Join(tmp, "base.ext4")
	cfg.Firecracker.SandboxDirectory = filepath.Join(tmp, "vms")

	if err := os.WriteFile(cfg.Storage.BaseRootfs, []byte("base"), 0o644); err != nil {
		t.Fatal(err)
	}

	m := storage.New(cfg)
	rootfs := m.Rootfs("vm-0001")
	if err := os.MkdirAll(filepath.Dir(rootfs), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(rootfs, []byte("existing"), 0o644); err != nil {
		t.Fatal(err)
	}

	err := m.CloneBaseRootfs("vm-0001")
	if err == nil {
		t.Fatal("expected error when rootfs already exists")
	}
}

func TestCloneBaseRootfs(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("cp --reflink=auto requires Linux")
	}

	tmp := t.TempDir()
	cfg := testConfig()
	cfg.Storage.BaseRootfs = filepath.Join(tmp, "base.ext4")
	cfg.Firecracker.SandboxDirectory = filepath.Join(tmp, "vms")

	base := []byte("sparse-base-image")
	if err := os.WriteFile(cfg.Storage.BaseRootfs, base, 0o644); err != nil {
		t.Fatal(err)
	}

	m := storage.New(cfg)
	if err := m.CloneBaseRootfs("vm-0001"); err != nil {
		t.Fatal(err)
	}

	got, err := os.ReadFile(m.Rootfs("vm-0001"))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(base) {
		t.Fatalf("cloned content = %q, want %q", got, base)
	}
}
