package allocate_test

import (
	"testing"

	"github.com/dwrth/knaller/internal/allocate"
	"github.com/dwrth/knaller/internal/state"
	"github.com/oklog/ulid/v2"
)

func TestNextSlot(t *testing.T) {
	uniqueID := func() string { return ulid.Make().String() }
	tests := []struct {
		name     string
		existing []state.Sandbox
		want     int
		wantErr  bool
	}{

		{
			name: "take empty slot inbetween",
			existing: []state.Sandbox{
				{Name: "test1", ID: uniqueID(), Slot: 1},
				{Name: "test3", ID: uniqueID(), Slot: 3},
			},
			want:    2,
			wantErr: false,
		},
		{
			name: "take empty 0 slot inbetween",
			existing: []state.Sandbox{
				{Name: "test1", ID: uniqueID(), Slot: 1},
				{Name: "test2", ID: uniqueID(), Slot: 0},
				{Name: "test3", ID: uniqueID(), Slot: 3},
			},
			want:    2,
			wantErr: false,
		},
		{
			name:    "no existing slots",
			want:    1,
			wantErr: false,
		},
		{
			name:     "no available slots",
			existing: buildMaxSandboxes(),
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := allocate.NextSlot(tt.existing)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("NextSlot() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("NextSlot() succeeded unexpectedly")
			}
			if got != tt.want {
				t.Errorf("NextSlot() = %v, want %v", got, tt.want)
			}
		})
	}
}

func buildMaxSandboxes() []state.Sandbox {
	existing := make([]state.Sandbox, 255)
	for i := range existing {
		existing[i] = state.Sandbox{Name: "test", Slot: i + 1}
	}

	return existing
}
