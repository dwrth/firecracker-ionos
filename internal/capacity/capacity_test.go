package capacity

import (
	"testing"

	"github.com/dwrth/knaller/internal/state"
)

func TestCapacity_AvailableMemoryMiB(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		want int64
		c    Capacity
	}{
		{
			name: "happy path",
			want: 12288,
			c: Capacity{
				HostMemoryMiB:      16384,
				ReservedMemoryMiB:  1024,
				AllocatedMemoryMiB: 3072,
			},
		},
		{
			name: "negative available memory",
			want: 0,
			c: Capacity{
				HostMemoryMiB:      16384,
				ReservedMemoryMiB:  20480,
				AllocatedMemoryMiB: 1024,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.c.AvailableMemoryMiB()
			if got != tt.want {
				t.Errorf("AvailableMemoryMiB() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCapacity_MaximumAllocatedVCPUs(t *testing.T) {
	tests := []struct {
		name string
		want int
		c    Capacity
	}{
		{
			name: "happy path",
			want: 16,
			c: Capacity{
				HostCPUs:           8,
				CPUOvercommitRatio: 2.0,
			},
		},
		{
			name: "zero host CPUs",
			want: 0,
			c: Capacity{
				HostCPUs:           0,
				CPUOvercommitRatio: 2.0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.c.MaximumAllocatedVCPUs()
			if got != tt.want {
				t.Errorf("MaximumAllocatedVCPUs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCapacity_AvailableVCPUs(t *testing.T) {
	tests := []struct {
		name string
		want int
		c    Capacity
	}{
		{
			name: "happy path",
			want: 10,
			c: Capacity{
				HostCPUs:           8,
				CPUOvercommitRatio: 2.0,
				AllocatedVCPUs:     6,
			},
		},
		{
			name: "negative available vCPUs",
			want: 0,
			c: Capacity{
				HostCPUs:           2,
				CPUOvercommitRatio: 2.0,
				AllocatedVCPUs:     6,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.c.AvailableVCPUs()
			if got != tt.want {
				t.Errorf("AvailableVCPUs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCapacity_sumAllocated(t *testing.T) {
	vcpus, mem := sumAllocated([]state.VM{
		{VCPUs: 2, MemoryMiB: 1024},
		{VCPUs: 1, MemoryMiB: 512},
	})
	if vcpus != 3 || mem != 1536 {
		t.Errorf("sumAllocated() = %d, %d, want 3, 1536", vcpus, mem)
	}
}
