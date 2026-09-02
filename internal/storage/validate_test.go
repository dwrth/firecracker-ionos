package storage_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dwrth/knaller/internal/storage"
)

func TestValidateSandboxStorageMissing(t *testing.T) {
	tmp := t.TempDir()
	cfg := testConfig()
	cfg.Firecracker.SandboxDirectory = filepath.Join(tmp, "vms")

	m := storage.New(cfg)
	err := m.ValidateSandboxStorage("vm-0001")
	if err == nil {
		t.Fatal("expected validation error")
	}
	if !strings.Contains(err.Error(), "rootfs") ||
		!strings.Contains(err.Error(), "kernel") ||
		!strings.Contains(err.Error(), "config") {
		t.Fatalf("error should list missing files: %v", err)
	}
}

func TestValidateSandboxStorageOK(t *testing.T) {
	tmp := t.TempDir()
	cfg := testConfig()
	cfg.Firecracker.SandboxDirectory = filepath.Join(tmp, "vms")

	m := storage.New(cfg)
	dir := m.SandboxDir("vm-0001")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"rootfs.ext4", "vmlinux", "config.json"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	if err := m.ValidateSandboxStorage("vm-0001"); err != nil {
		t.Fatalf("ValidateSandboxStorage() = %v", err)
	}
}
