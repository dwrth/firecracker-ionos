package config

import (
	"errors"
	"net/netip"
	"os"
)

func (c *Config) Validate() error {
	if c == nil {
		return errors.New("config: config is nil")
	}

	if c.Storage.BaseRootfs == "" {
		return errors.New("config: storage.base_rootfs is required")
	}
	if _, err := os.Stat(c.Storage.BaseRootfs); err != nil {
		return errors.New("config: storage.base_rootfs does not exist")
	}
	if c.Storage.Kernel == "" {
		return errors.New("config: storage.kernel is required")
	}
	if _, err := os.Stat(c.Storage.Kernel); err != nil {
		return errors.New("config: storage.kernel does not exist")
	}

	if c.Scheduler.CPUOvercommitRatio <= 0 {
		return errors.New("config: scheduler.cpu_overcommit_ratio must be greater than 0")
	}
	if c.Scheduler.HostMemoryReserveMib <= 0 {
		return errors.New("config: scheduler.host_memory_reserve_mib must be greater than 0")
	}
	if c.Scheduler.MinimumVMMemoryMib <= 0 {
		return errors.New("config: scheduler.minimum_vm_memory_mib must be greater than 0")
	}
	if c.Scheduler.MaximumVMMemoryMib <= 0 {
		return errors.New("config: scheduler.maximum_vm_memory_mib must be greater than 0")
	}
	if c.Scheduler.MaximumVMMemoryMib < c.Scheduler.MinimumVMMemoryMib {
		return errors.New("config: scheduler.maximum_vm_memory_mib must be greater than or equal to scheduler.minimum_vm_memory_mib")
	}
	if c.Scheduler.MaximumVMVcpus <= 0 {
		return errors.New("config: scheduler.maximum_vm_vcpus must be greater than 0")
	}

	if _, err := netip.ParsePrefix(c.Network.GuestCidr); err != nil {
		return errors.New("config: network.guest_cidr is invalid")
	}
	if _, err := netip.ParsePrefix(c.Network.TransitCidr); err != nil {
		return errors.New("config: network.transit_cidr is invalid")
	}

	if c.Jailer.BaseDir == "" {
		return errors.New("config: jailer.base_dir is required")
	}
	if c.Jailer.UidStart <= 1000 {
		return errors.New("config: jailer.uid_start must be greater than 1000")
	}
	if c.Jailer.GidStart <= 1000 {
		return errors.New("config: jailer.gid_start must be greater than 1000")
	}

	if c.State.Directory == "" {
		return errors.New("config: state.directory is required")
	}

	if c.Firecracker.VmDirectory == "" {
		return errors.New("config: firecracker.vm_directory is required")
	}

	return nil
}
