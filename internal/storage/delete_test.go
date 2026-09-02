package storage_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/dwrth/knaller/internal/storage"
)

func TestDeleteSandboxStorage(t *testing.T) {
	tmp := t.TempDir()
	cfg := testConfig()
	cfg.Firecracker.SandboxDirectory = filepath.Join(tmp, "vms")

	m := storage.New(cfg)
	dir := m.SandboxDir("vm-0001")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "rootfs.ext4"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := m.DeleteSandboxStorage("vm-0001"); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		t.Fatalf("sandbox dir still exists after delete: %v", err)
	}
}

func TestDeleteSandboxStorageMissing(t *testing.T) {
	tmp := t.TempDir()
	cfg := testConfig()
	cfg.Firecracker.SandboxDirectory = filepath.Join(tmp, "vms")

	m := storage.New(cfg)
	if err := m.DeleteSandboxStorage("vm-0001"); err != nil {
		t.Fatalf("DeleteSandboxStorage() on missing dir = %v", err)
	}
}
