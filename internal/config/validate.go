package config

import (
	"errors"
	"net/netip"
	"os"
)

// Validate checks that required fields are set and paths exist.
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
	if c.Scheduler.HostMemoryReserveMiB <= 0 {
		return errors.New("config: scheduler.host_memory_reserve_mib must be greater than 0")
	}
	if c.Scheduler.MinimumSandboxMemoryMiB <= 0 {
		return errors.New("config: scheduler.minimum_sandbox_memory_mib must be greater than 0")
	}
	if c.Scheduler.MaximumSandboxMemoryMiB <= 0 {
		return errors.New("config: scheduler.maximum_sandbox_memory_mib must be greater than 0")
	}
	if c.Scheduler.MaximumSandboxMemoryMiB < c.Scheduler.MinimumSandboxMemoryMiB {
		return errors.New("config: scheduler.maximum_sandbox_memory_mib must be greater than or equal to scheduler.minimum_sandbox_memory_mib")
	}
	if c.Scheduler.MaximumSandboxVCPUs <= 0 {
		return errors.New("config: scheduler.maximum_sandbox_vcpus must be greater than 0")
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

	if c.Firecracker.SandboxDirectory == "" {
		return errors.New("config: firecracker.sandbox_directory is required")
	}

	return nil
}
