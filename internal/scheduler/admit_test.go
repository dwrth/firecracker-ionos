package scheduler_test

import (
	"errors"
	"testing"

	"github.com/dwrth/knaller/internal/capacity"
	"github.com/dwrth/knaller/internal/config"
	"github.com/dwrth/knaller/internal/scheduler"
)

func testCfg() *config.Config {
	return &config.Config{
		Scheduler: config.Scheduler{
			MinimumVMMemoryMiB: 128,
			MaximumVMMemoryMiB: 2048,
			MaximumVMVCPUs:     4,
		},
	}
}

func TestAdmit(t *testing.T) {
	cfg := testCfg()
	cap := capacity.Capacity{}

	tests := []struct {
		name    string
		req     scheduler.Request
		wantErr error
	}{
		{"ok", scheduler.Request{VCPUs: 2, MemoryMiB: 512}, nil},
		{"zero vcpus", scheduler.Request{VCPUs: 0, MemoryMiB: 512}, scheduler.ErrInvalidVCPUs},
		{"too many vcpus", scheduler.Request{VCPUs: 8, MemoryMiB: 512}, scheduler.ErrTooManyVCPUs},
		{"memory too small", scheduler.Request{VCPUs: 1, MemoryMiB: 64}, scheduler.ErrMemoryTooSmall},
		{"memory too large", scheduler.Request{VCPUs: 1, MemoryMiB: 4096}, scheduler.ErrMemoryTooLarge},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := scheduler.Admit(cfg, cap, tt.req)
			if tt.wantErr == nil {
				if err != nil {
					t.Errorf("Admit() = %v, want nil", err)
				}
				return
			}
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Admit() = %v, want %v", err, tt.wantErr)
			}
		})
	}
}
