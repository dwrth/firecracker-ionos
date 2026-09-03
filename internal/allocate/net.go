package allocate

import (
	"fmt"
	"net/netip"
)

// Net holds a carved /30 subnet and host and peer addresses.
type Net struct {
	Subnet string
	HostIP string
	PeerIP string
}

// GuestNet carves a /30 guest subnet from guestCIDR for id.
func GuestNet(slot int, guestCIDR string) (Net, error) {
	return subnetFor(slot, guestCIDR)
}

// TransitNet carves a /30 transit subnet from transitCIDR for id.
func TransitNet(slot int, transitCIDR string) (Net, error) {
	return subnetFor(slot, transitCIDR)
}

func subnetFor(slot int, parentCIDR string) (Net, error) {
	parent, err := netip.ParsePrefix(parentCIDR)
	if err != nil {
		return Net{}, fmt.Errorf("allocate: invalid cidr %q: %w", parentCIDR, err)
	}
	addr := parent.Addr()
	if !addr.Is4() {
		return Net{}, fmt.Errorf("allocate: only IPv4 cidrs supported, got %q", parentCIDR)
	}

	octets := addr.As4()
	octets[2] = byte(slot)
	octets[3] = 0
	network := netip.AddrFrom4(octets)
	subnet := netip.PrefixFrom(network, 30)
	if !parent.Contains(network) {
		return Net{}, fmt.Errorf("allocate: %s nit inside %s", subnet, parent)
	}

	octets[3] = 1
	host := netip.AddrFrom4(octets)
	octets[3] = 2
	peer := netip.AddrFrom4(octets)

	return Net{
		Subnet: subnet.String(),
		HostIP: host.String(),
		PeerIP: peer.String(),
	}, nil
}
