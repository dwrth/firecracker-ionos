package allocate_test

import (
	"testing"

	"github.com/dwrth/knaller/internal/allocate"
)

func TestGuestNet(t *testing.T) {
	got, err := allocate.GuestNet("vm-0001", "172.16.0.0/16")
	if err != nil {
		t.Fatal(err)
	}
	if got.Subnet != "172.16.1.0/30" || got.HostIP != "172.16.1.1" || got.PeerIP != "172.16.1.2" {
		t.Errorf("GuestNet() = %v, want %v", got, allocate.Net{Subnet: "172.16.1.0/30", HostIP: "172.16.1.1", PeerIP: "172.16.1.2"})
	}
}

func TestTransitNet(t *testing.T) {
	got, err := allocate.TransitNet("vm-0002", "10.200.0.0/16")
	if err != nil {
		t.Fatal(err)
	}
	if got.Subnet != "10.200.2.0/30" || got.HostIP != "10.200.2.1" || got.PeerIP != "10.200.2.2" {
		t.Errorf("TransitNet() = %v, want %v", got, allocate.Net{Subnet: "10.200.2.0/30", HostIP: "10.200.2.1", PeerIP: "10.200.2.2"})
	}
}
