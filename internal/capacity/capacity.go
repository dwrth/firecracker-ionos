// Package capacity will report host and VM resource capacity.
package capacity

import (
	"fmt"
	"io"
	"math"
	"os"

	"github.com/dwrth/knaller/internal/config"
	"github.com/dwrth/knaller/internal/state"
)

type Capacity struct {
	HostCPUs           int
	HostMemoryMiB      int64
	ReservedMemoryMiB  int64
	AllocatedVCPUs     int
	AllocatedMemoryMiB int64
	CPUOvercommitRatio float64
}

func (c *Capacity) AvailableMemoryMiB() int64 {
	return int64(math.Max(0, float64(c.HostMemoryMiB-c.ReservedMemoryMiB-c.AllocatedMemoryMiB)))
}

func (c *Capacity) MaximumAllocatedVCPUs() int {
	return int(float64(c.HostCPUs) * c.CPUOvercommitRatio)
}

func (c *Capacity) AvailableVCPUs() int {
	return int(math.Max(0, float64(c.MaximumAllocatedVCPUs()-c.AllocatedVCPUs)))
}

func Collect(cfg *config.Config) (Capacity, error) {
	meminfo, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return Capacity{}, fmt.Errorf("capacity: failed to read meminfo: %v", err)
	}

	memtotal, err := parseMemTotalMiB(meminfo)
	if err != nil {
		return Capacity{}, fmt.Errorf("capacity: failed to parse meminfo: %v", err)
	}

	vms, err := state.New(cfg.State.Directory).List()
	if err != nil {
		return Capacity{}, fmt.Errorf("capacity: failed to list VMs: %w", err)
	}

	allocatedVCPUs, allocatedMemoryMiB := sumAllocated(vms)

	capacity := Capacity{
		HostCPUs:           hostCPUs(),
		HostMemoryMiB:      memtotal,
		ReservedMemoryMiB:  int64(cfg.Scheduler.HostMemoryReserveMib),
		CPUOvercommitRatio: cfg.Scheduler.CPUOvercommitRatio,
		AllocatedVCPUs:     allocatedVCPUs,
		AllocatedMemoryMiB: allocatedMemoryMiB,
	}

	return capacity, nil
}

func Print(w io.Writer, capacity Capacity) {
	fmt.Fprintf(w, "===== Host Capacity =====\n")
	fmt.Fprintf(w, " Host CPUs            : %d\n", capacity.HostCPUs)
	fmt.Fprintf(w, " Host memory          : %d MiB\n", capacity.HostMemoryMiB)
	fmt.Fprintf(w, " Reserved memory      : %d MiB\n", capacity.ReservedMemoryMiB)
	fmt.Fprintf(w, " CPU overcommit ratio : %.2f\n", capacity.CPUOvercommitRatio)
	fmt.Fprintf(w, "-----------------------------\n")
	fmt.Fprintf(w, " Allocated vCPUs      : %d\n", capacity.AllocatedVCPUs)
	fmt.Fprintf(w, " Allocated memory     : %d MiB\n", capacity.AllocatedMemoryMiB)
	fmt.Fprintf(w, " Maximum vCPUs        : %d\n", capacity.MaximumAllocatedVCPUs())
	fmt.Fprintf(w, " Available vCPUs      : %d\n", capacity.AvailableVCPUs())
	fmt.Fprintf(w, " Available memory     : %d MiB\n", capacity.AvailableMemoryMiB())
	fmt.Fprintf(w, "=============================\n")
}

func sumAllocated(vms []state.VM) (vcpus int, memoryMiB int64) {
	for _, vm := range vms {
		vcpus += vm.VCPUs
		memoryMiB += int64(vm.MemoryMiB)
	}
	return vcpus, memoryMiB
}
