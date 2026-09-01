package allocate_test

import (
	"testing"

	"github.com/dwrth/knaller/internal/allocate"
	"github.com/dwrth/knaller/internal/config"
	"github.com/dwrth/knaller/internal/state"
)

func TestBuild(t *testing.T) {
	cfg := &config.Config{
		Jailer:  config.Jailer{UidStart: 12000, GidStart: 12000},
		Network: config.Network{GuestCidr: "172.16.0.0/16", TransitCidr: "10.200.0.0/16"},
	}
	vm, err := allocate.Build(nil, cfg, "web-1", 2, 512)
	if err != nil {
		t.Fatal(err)
	}
	if vm.ID != "vm-0001" || vm.UID != 12001 || vm.GuestIP != "172.16.1.2" {
		t.Fatalf("Build() = %+v", vm)
	}
	if vm.Status != state.StatusStopped {
		t.Fatalf("status = %v, want stopped", vm.Status)
	}
}
