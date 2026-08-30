package allocate

import (
	"fmt"
	"net/netip"
)

type Net struct {
	Subnet string
	HostIP string
	PeerIP string
}

func GuestNet(id, guestCIDR string) (Net, error) {
	return subnetFor(id, guestCIDR)
}

func TransitNet(id, transitCIDR string) (Net, error) {
	return subnetFor(id, transitCIDR)
}

func subnetFor(id, parentCIDR string) (Net, error) {
	n, err := parseID(id)
	if err != nil {
		return Net{}, err
	}
	if n < 1 || n > 255 {
		return Net{}, fmt.Errorf("allocate: vm index %d out of range for /30 carving", n)
	}

	parent, err := netip.ParsePrefix(parentCIDR)
	if err != nil {
		return Net{}, fmt.Errorf("allocate: invalid cidr %q: %w", parentCIDR, err)
	}
	addr := parent.Addr()
	if !addr.Is4() {
		return Net{}, fmt.Errorf("allocate: only IPv4 cidrs supported, got %q", parentCIDR)
	}

	octets := addr.As4()
	octets[2] = byte(n)
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
