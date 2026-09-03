package allocate

import (
	"fmt"

	"github.com/dwrth/knaller/internal/config"
	"github.com/dwrth/knaller/internal/state"
)

// Build allocates a sandbox ID, slot, credentials, networking, and interface names.
func Build(existing []state.Sandbox, cfg *config.Config, name string, vcpus, memoryMiB int) (state.Sandbox, error) {
	if name == "" {
		return state.Sandbox{}, fmt.Errorf("allocate: name is required")
	}

	id, err := NewSandboxID(existing)
	if err != nil {
		return state.Sandbox{}, err
	}
	slot, err := NextSlot(existing)
	if err != nil {
		return state.Sandbox{}, err
	}

	uid, gid := Creds(slot, cfg.Jailer.UidStart, cfg.Jailer.GidStart)
	guest, err := GuestNet(slot, cfg.Network.GuestCidr)
	if err != nil {
		return state.Sandbox{}, err
	}
	transit, err := TransitNet(slot, cfg.Network.TransitCidr)
	if err != nil {
		return state.Sandbox{}, err
	}
	names := InterfaceNames(slot)

	return state.Sandbox{
		ID:            id,
		Name:          name,
		Slot:          slot,
		UID:           uid,
		GID:           gid,
		VCPUs:         vcpus,
		MemoryMiB:     memoryMiB,
		GuestIP:       guest.PeerIP,
		GatewayIP:     guest.HostIP,
		GuestSubnet:   guest.Subnet,
		TransitHostIP: transit.HostIP,
		TransitNSIP:   transit.PeerIP,
		TransitSubnet: transit.Subnet,
		Namespace:     names.Namespace,
		HostVeth:      names.HostVeth,
		NSVeth:        names.NSVeth,
		TAP:           names.TAP,
		DesiredState:  state.DesiredStopped,
		ObservedState: state.ObservedRequested,
	}, nil

}
