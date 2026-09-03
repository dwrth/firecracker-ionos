package allocate

import (
	"fmt"

	"github.com/dwrth/knaller/internal/state"
)

const maxSlot = 255

// NextSlot finds and returns the lowest unused slot number among the provided sandboxes.
//
// Each sandbox's Slot is a node-local integer resource (starting from 1) used for slot allocation.
// The returned integer is guaranteed not to collide with any Slot currently present in 'existing'.
func NextSlot(existing []state.Sandbox) (int, error) {
	used := make(map[int]struct{}, len(existing))
	for _, sandbox := range existing {
		if sandbox.Slot >= 1 {
			used[sandbox.Slot] = struct{}{}
		}
	}

	for n := 1; ; n++ {
		if n > maxSlot {
			return 0, fmt.Errorf("allocate: no slots available (max %d)", maxSlot)
		}
		if _, ok := used[n]; !ok {
			return n, nil
		}
	}
}
