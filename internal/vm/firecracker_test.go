package vm_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/dwrth/knaller/internal/state"
	"github.com/dwrth/knaller/internal/vm"
)

func TestGenerateFirecrackerConfig_matchesGolden(t *testing.T) {
	tests := []struct {
		name     string
		sandbox  state.Sandbox
		guestDir string
	}{
		{
			name: "vm1",
			sandbox: state.Sandbox{
				VCPUs:     1,
				MemoryMiB: 256,
				GuestIP:   "172.16.1.2",
				TAP:       "tap0",
			},
			guestDir: "vm1",
		},
		{
			name: "vm2",
			sandbox: state.Sandbox{
				VCPUs:     1,
				MemoryMiB: 256,
				GuestIP:   "172.16.2.2",
				TAP:       "tap0",
			},
			guestDir: "vm2",
		},
		{
			name: "plan vm-0004",
			sandbox: state.Sandbox{
				ID:        "vm-0004",
				Name:      "worker-1",
				VCPUs:     2,
				MemoryMiB: 1024,
				GuestIP:   "172.16.4.2",
				TAP:       "tap0",
			},
			guestDir: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := vm.GenerateFirecrackerConfig(tt.sandbox)
			if err != nil {
				t.Fatal(err)
			}

			if tt.guestDir != "" {
				want := loadGoldenConfig(t, tt.guestDir)
				if !reflect.DeepEqual(got, want) {
					t.Fatalf("GenerateFirecrackerConfig() = %+v, want %+v", got, want)
				}
				return
			}

			if got.MachineConfig.VCPUCount != 2 || got.MachineConfig.MemSizeMiB != 1024 {
				t.Fatalf("machine-config = %+v, want 2 vcpus and 1024 MiB", got.MachineConfig)
			}
			if got.NetworkInterfaces[0].GuestMac != "06:00:AC:10:04:02" {
				t.Fatalf("guest_mac = %q, want 06:00:AC:10:04:02", got.NetworkInterfaces[0].GuestMac)
			}
		})
	}
}

func TestGenerateFirecrackerConfig_invalidGuestIP(t *testing.T) {
	_, err := vm.GenerateFirecrackerConfig(state.Sandbox{
		VCPUs:     1,
		MemoryMiB: 256,
		GuestIP:   "not-an-ip",
		TAP:       "tap0",
	})
	if err == nil {
		t.Fatal("expected error for invalid guest ip")
	}
}

func loadGoldenConfig(t *testing.T, guest string) vm.FirecrackerConfig {
	t.Helper()

	path := filepath.Join("..", "..", "guests", guest, "config.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	var cfg vm.FirecrackerConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		t.Fatal(err)
	}
	return cfg
}
