// Package state will manage the state of the system.
package state

// Status is the runtime state of a VM.
type Status string

const (
	StatusStopped Status = "stopped"
	StatusRunning Status = "running"
)

// VM describes a microVM and its allocated resources.
type VM struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	UID           int    `json:"uid"`
	GID           int    `json:"gid"`
	VCPUs         int    `json:"vcpus"`
	MemoryMiB     int    `json:"memory_mib"`
	GuestIP       string `json:"guest_ip"`
	GatewayIP     string `json:"gateway_ip"`
	GuestSubnet   string `json:"guest_subnet"`
	TransitHostIP string `json:"transit_host_ip"`
	TransitNSIP   string `json:"transit_ns_ip"`
	TransitSubnet string `json:"transit_subnet"`
	Namespace     string `json:"namespace"`
	HostVeth      string `json:"host_veth"`
	NSVeth        string `json:"ns_veth"`
	TAP           string `json:"tap"`
	Status        Status `json:"status"`
}
