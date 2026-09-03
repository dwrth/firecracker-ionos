package allocate

import (
	"fmt"

	"github.com/dwrth/knaller/internal/config"
	"github.com/dwrth/knaller/internal/state"
)

// Build allocates ids, credentials, networking, and interface names for a new sandbox.
func Build(existing []state.Sandbox, cfg *config.Config, name string, vcpus, memoryMiB int) (state.Sandbox, error) {
	if name == "" {
		return state.Sandbox{}, fmt.Errorf("allocate: name is required")
	}

	id, err := NextSlotKey(existing)
	if err != nil {
		return state.Sandbox{}, err
	}
	uid, gid, err := Creds(id, cfg.Jailer.UidStart, cfg.Jailer.GidStart)
	if err != nil {
		return state.Sandbox{}, err
	}
	guest, err := GuestNet(id, cfg.Network.GuestCidr)
	if err != nil {
		return state.Sandbox{}, err
	}
	transit, err := TransitNet(id, cfg.Network.TransitCidr)
	if err != nil {
		return state.Sandbox{}, err
	}
	names, err := InterfaceNames(id)
	if err != nil {
		return state.Sandbox{}, err
	}

	return state.Sandbox{
		ID:            id,
		Name:          name,
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
