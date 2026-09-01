// Package config will manage the configuration of the system.
package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Storage     Storage     `yaml:"storage"`
	Scheduler   Scheduler   `yaml:"scheduler"`
	Network     Network     `yaml:"network"`
	Jailer      Jailer      `yaml:"jailer"`
	State       State       `yaml:"state"`
	Firecracker Firecracker `yaml:"firecracker"`
}

type Storage struct {
	BaseRootfs string `yaml:"base_rootfs"`
	Kernel     string `yaml:"kernel"`
}

type Scheduler struct {
	HostMemoryReserveMiB int     `yaml:"host_memory_reserve_mib"`
	CPUOvercommitRatio   float64 `yaml:"cpu_overcommit_ratio"`
	MinimumVMMemoryMiB   int     `yaml:"minimum_vm_memory_mib"`
	MaximumVMMemoryMiB   int     `yaml:"maximum_vm_memory_mib"`
	MaximumVMVCPUs       int     `yaml:"maximum_vm_vcpus"`
}

type Network struct {
	GuestCidr   string `yaml:"guest_cidr"`
	TransitCidr string `yaml:"transit_cidr"`
}

type Jailer struct {
	BaseDir  string `yaml:"base_dir"`
	UidStart int    `yaml:"uid_start"`
	GidStart int    `yaml:"gid_start"`
}

type State struct {
	Directory string `yaml:"directory"`
}

type Firecracker struct {
	VmDirectory string `yaml:"vm_directory"`
}

// Load reads and validates configuration from path.
func (c *Config) Load(path string) error {
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		return err
	}

	return c.Validate()
}
