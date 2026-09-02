package storage_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dwrth/knaller/internal/state"
	"github.com/dwrth/knaller/internal/storage"
)

func testSandbox() state.Sandbox {
	return state.Sandbox{
		ID:        "vm-0004",
		Name:      "worker-1",
		GuestIP:   "172.16.4.2",
		GatewayIP: "172.16.4.1",
	}
}

func TestCustomizeMountedRootfs(t *testing.T) {
	root := t.TempDir()

	if err := storage.CustomizeMountedRootfsForTest(root, testSandbox()); err != nil {
		t.Fatal(err)
	}

	hostname, err := os.ReadFile(filepath.Join(root, "etc/hostname"))
	if err != nil {
		t.Fatal(err)
	}
	if got, want := strings.TrimSpace(string(hostname)), "vm-0004"; got != want {
		t.Errorf("hostname = %q, want %q", got, want)
	}

	hosts, err := os.ReadFile(filepath.Join(root, "etc/hosts"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(hosts), "127.0.1.1 vm-0004") {
		t.Errorf("hosts missing sandbox entry: %s", hosts)
	}

	network, err := os.ReadFile(filepath.Join(root, "etc/systemd/network/10-eth0.network"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(network), "Address=172.16.4.2/30") {
		t.Errorf("network missing guest address: %s", network)
	}
	if !strings.Contains(string(network), "Gateway=172.16.4.1") {
		t.Errorf("network missing gateway: %s", network)
	}

	machineID, err := os.ReadFile(filepath.Join(root, "etc/machine-id"))
	if err != nil {
		t.Fatal(err)
	}
	id := strings.TrimSpace(string(machineID))
	if len(id) != 32 {
		t.Errorf("machine-id length = %d, want 32", len(id))
	}
}

func TestCustomizeMountedRootfsUniqueMachineID(t *testing.T) {
	root := t.TempDir()
	sandbox := testSandbox()

	if err := storage.CustomizeMountedRootfsForTest(root, sandbox); err != nil {
		t.Fatal(err)
	}
	first, err := os.ReadFile(filepath.Join(root, "etc/machine-id"))
	if err != nil {
		t.Fatal(err)
	}

	if err := storage.CustomizeMountedRootfsForTest(root, sandbox); err != nil {
		t.Fatal(err)
	}
	second, err := os.ReadFile(filepath.Join(root, "etc/machine-id"))
	if err != nil {
		t.Fatal(err)
	}

	if string(first) == string(second) {
		t.Fatal("expected different machine-id on second customize")
	}
}
