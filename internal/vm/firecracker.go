// Package vm builds per-VM Firecracker configuration from allocated state.
package vm

import (
	"fmt"
	"net/netip"

	"github.com/dwrth/knaller/internal/state"
)

type FirecrackerConfig struct {
	BootSource        BootSource         `json:"boot-source"`
	Drives            []Drive            `json:"drives"`
	MachineConfig     MachineConfig      `json:"machine-config"`
	NetworkInterfaces []NetworkInterface `json:"network-interfaces"`
}

type BootSource struct {
	KernelImagePath string `json:"kernel_image_path"`
	BootArgs        string `json:"boot_args"`
}

type Drive struct {
	DriveID      string `json:"drive_id"`
	PathOnHost   string `json:"path_on_host"`
	IsRootDevice bool   `json:"is_root_device"`
	IsReadOnly   bool   `json:"is_read_only"`
}

type MachineConfig struct {
	VCPUCount  int  `json:"vcpu_count"`
	MemSizeMiB int  `json:"mem_size_mib"`
	SMT        bool `json:"smt"`
}

type NetworkInterface struct {
	IFaceID     string `json:"iface_id"`
	GuestMac    string `json:"guest_mac"`
	HostDevName string `json:"host_dev_name"`
}

func GenerateFirecrackerConfig(vm state.VM) (FirecrackerConfig, error) {
	mac, err := guestMAC(vm.GuestIP)
	if err != nil {
		return FirecrackerConfig{}, err
	}
	return FirecrackerConfig{
		BootSource: BootSource{
			KernelImagePath: "/vmlinux",
			BootArgs:        "console=ttyS0 reboot=k panic=1",
		},
		Drives: []Drive{{
			DriveID:      "rootfs",
			PathOnHost:   "/rootfs.ext4",
			IsRootDevice: true,
			IsReadOnly:   false,
		}},
		MachineConfig: MachineConfig{
			VCPUCount:  vm.VCPUs,
			MemSizeMiB: vm.MemoryMiB,
			SMT:        false,
		},
		NetworkInterfaces: []NetworkInterface{{
			IFaceID:     "eth0",
			GuestMac:    mac,
			HostDevName: vm.TAP,
		}},
	}, nil
}

func guestMAC(ip string) (string, error) {
	addr, err := netip.ParseAddr(ip)
	if err != nil {
		return "", fmt.Errorf("vm: guest ip: %w", err)
	}
	if !addr.Is4() {
		return "", fmt.Errorf("vm: guest ip %q is not IPv4", ip)
	}

	o := addr.As4()
	return fmt.Sprintf("06:00:%02X:%02X:%02X:%02X", o[0], o[1], o[2], o[3]), nil
}
