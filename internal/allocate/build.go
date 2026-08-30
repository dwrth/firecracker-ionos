package allocate

import (
	"fmt"

	"github.com/dwrth/knaller/internal/config"
	"github.com/dwrth/knaller/internal/state"
)

func Build(existing []state.VM, cfg *config.Config, name string, vcpus, memoryMiB int) (state.VM, error) {
	if name == "" {
		return state.VM{}, fmt.Errorf("allocate: name is required")
	}

	id, err := NextID(existing)
	if err != nil {
		return state.VM{}, err
	}
	uid, gid, err := Creds(id, cfg.Jailer.UidStart, cfg.Jailer.GidStart)
	if err != nil {
		return state.VM{}, err
	}
	guest, err := GuestNet(id, cfg.Network.GuestCidr)
	if err != nil {
		return state.VM{}, err
	}
	transit, err := TransitNet(id, cfg.Network.TransitCidr)
	if err != nil {
		return state.VM{}, err
	}
	names, err := InterfaceNames(id)
	if err != nil {
		return state.VM{}, err
	}

	return state.VM{
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
		Status:        state.StatusStopped,
	}, nil

}
