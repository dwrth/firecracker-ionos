package allocate_test

import (
	"testing"

	"github.com/dwrth/knaller/internal/allocate"
	"github.com/dwrth/knaller/internal/state"
)

func TestNextID(t *testing.T) {
	tests := []struct {
		name     string
		existing []state.VM
		want     string
	}{
		{"empty", nil, "vm-0001"},
		{"sequential", []state.VM{{ID: "vm-0001"}, {ID: "vm-0002"}}, "vm-0003"},
		{"hole", []state.VM{{ID: "vm-0001"}, {ID: "vm-0003"}}, "vm-0002"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := allocate.NextID(tt.existing)
			if err != nil {
				t.Fatal(err)
			}
			if got != tt.want {
				t.Errorf("NextID() = %q, want %q", got, tt.want)
			}
		})
	}
}
