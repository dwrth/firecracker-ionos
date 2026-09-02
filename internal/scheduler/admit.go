// Package scheduler will schedule sandboxes.
package scheduler

import (
	"errors"
	"fmt"

	"github.com/dwrth/knaller/internal/capacity"
	"github.com/dwrth/knaller/internal/config"
)

// Request describes resources requested for a new sandbox.
type Request struct {
	VCPUs     int
	MemoryMiB int
}

var (
	ErrInvalidVCPUs       = errors.New("scheduler: vcpus must be at least 1")
	ErrTooManyVCPUs       = errors.New("scheduler: vcpus exceed maximum_sandbox_vcpus")
	ErrMemoryTooSmall     = errors.New("scheduler: memory below minimum_sandbox_memory_mib")
	ErrMemoryTooLarge     = errors.New("scheduler: memory above maximum_sandbox_memory_mib")
	ErrInsufficientMemory = errors.New("scheduler: insufficient memory")
	ErrInsufficientVCPUs  = errors.New("scheduler: insufficient vcpus")
)

// Admit checks whether req fits configured limits and available capacity.
func Admit(cfg *config.Config, cap capacity.Capacity, req Request) error {
	if req.VCPUs < 1 {
		return ErrInvalidVCPUs
	}
	if req.VCPUs > cfg.Scheduler.MaximumSandboxVCPUs {
		return fmt.Errorf("%w: requested %d, max %d", ErrTooManyVCPUs, req.VCPUs, cfg.Scheduler.MaximumSandboxVCPUs)
	}
	if req.MemoryMiB < cfg.Scheduler.MinimumSandboxMemoryMiB {
		return fmt.Errorf("%w: requested %d, min %d", ErrMemoryTooSmall, req.MemoryMiB, cfg.Scheduler.MinimumSandboxMemoryMiB)
	}
	if req.MemoryMiB > cfg.Scheduler.MaximumSandboxMemoryMiB {
		return fmt.Errorf("%w: requested %d, max %d", ErrMemoryTooLarge, req.MemoryMiB, cfg.Scheduler.MaximumSandboxMemoryMiB)
	}
	if int64(req.MemoryMiB) > cap.AvailableMemoryMiB() {
		return fmt.Errorf("%w: requested %d MiB, available %d MiB", ErrInsufficientMemory, req.MemoryMiB, cap.AvailableMemoryMiB())
	}
	if req.VCPUs > cap.AvailableVCPUs() {
		return fmt.Errorf("%w: requested %d, available %d", ErrInsufficientVCPUs, req.VCPUs, cap.AvailableVCPUs())
	}
	return nil
}
