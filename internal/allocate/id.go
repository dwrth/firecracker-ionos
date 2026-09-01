// Package allocate will allocate VMs.
package allocate

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/dwrth/knaller/internal/state"
)

// NextID returns the lowest unused vm-NNNN id not present in existing.
func NextID(existing []state.VM) (string, error) {
	used := make(map[int]struct{}, len(existing))
	for _, vm := range existing {
		n, err := parseID(vm.ID)
		if err != nil {
			return "", err
		}
		used[n] = struct{}{}
	}

	for n := 1; ; n++ {
		if _, ok := used[n]; !ok {
			return formatID(n), nil
		}
	}
}

func parseID(id string) (int, error) {
	if !strings.HasPrefix(id, "vm-") {
		return 0, fmt.Errorf("allocate: invalid vm id: %q", id)
	}
	n, err := strconv.Atoi(id[3:])
	if err != nil || n < 0 {
		return 0, fmt.Errorf("allocate: invalid vm id: %q", id)
	}
	return n, nil
}

func formatID(n int) string {
	return fmt.Sprintf("vm-%04d", n)
}
