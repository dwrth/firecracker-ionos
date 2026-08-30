// Package scheduler will schedule VMs.
package scheduler

import (
	"errors"
	"fmt"

	"github.com/dwrth/knaller/internal/capacity"
	"github.com/dwrth/knaller/internal/config"
)

type Request struct {
	VCPUs     int
	MemoryMiB int
}

var (
	ErrInvalidVCPUs   = errors.New("scheduler: vcpus must be at least 1")
	ErrTooManyVCPUs   = errors.New("scheduler: vcpus exceed maximum_vm_vcpus")
	ErrMemoryTooSmall = errors.New("scheduler: memory below minimum_vm_memory_mib")
	ErrMemoryTooLarge = errors.New("scheduler: memory above maximum_vm_memory_mib")
)

func Admit(cfg *config.Config, cap capacity.Capacity, req Request) error {
	if req.VCPUs < 1 {
		return ErrInvalidVCPUs
	}
	if req.VCPUs > cfg.Scheduler.MaximumVMVCPUs {
		return fmt.Errorf("%w: requested %d, max %d", ErrTooManyVCPUs, req.VCPUs, cfg.Scheduler.MaximumVMVCPUs)
	}
	if req.MemoryMiB < cfg.Scheduler.MinimumVMMemoryMiB {
		return fmt.Errorf("%w: requested %d, min %d", ErrMemoryTooSmall, req.MemoryMiB, cfg.Scheduler.MinimumVMMemoryMiB)
	}
	if req.MemoryMiB > cfg.Scheduler.MaximumVMMemoryMiB {
		return fmt.Errorf("%w: requested %d, max %d", ErrMemoryTooLarge, req.MemoryMiB, cfg.Scheduler.MaximumVMMemoryMiB)
	}

	_ = cap
	return nil
}
