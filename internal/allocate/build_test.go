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
	sandbox, err := allocate.Build(nil, cfg, "web-1", 2, 512)
	if err != nil {
		t.Fatal(err)
	}
	if sandbox.ID != "vm-0001" || sandbox.UID != 12001 || sandbox.GuestIP != "172.16.1.2" {
		t.Fatalf("Build() = %+v", sandbox)
	}
	if sandbox.DesiredState != state.DesiredStopped {
		t.Fatalf("desired state = %v, want stopped", sandbox.DesiredState)
	}
	if sandbox.ObservedState != state.ObservedRequested {
		t.Fatalf("observed state = %v, want requested", sandbox.ObservedState)
	}
}
