// Package allocate provides utilities for generating and managing unique sandbox identifiers.
package allocate

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/dwrth/knaller/internal/state"
	"github.com/oklog/ulid/v2"
)

// NextSlotKey returns the lowest unused vm-NNNN id not present in existing.
func NextSlotKey(existing []state.Sandbox) (string, error) {
	used := make(map[int]struct{}, len(existing))
	for _, sandbox := range existing {
		n, err := parseSlotKey(sandbox.ID)
		if err != nil {
			return "", err
		}
		used[n] = struct{}{}
	}

	for n := 1; ; n++ {
		if _, ok := used[n]; !ok {
			return formatSlotKey(n), nil
		}
	}
}

// NewSandboxID returns a globally unique sandbox ID not present in existing.
func NewSandboxID(existing []state.Sandbox) (string, error) {
	used := make(map[string]struct{}, len(existing))
	for _, sb := range existing {
		used[sb.ID] = struct{}{}
	}

	for range 10 {
		id := ulid.Make().String()
		if _, ok := used[id]; !ok {
			return id, nil
		}
	}
	return "", fmt.Errorf("allocate: failed to generate unique sandbox id")
}

func parseSlotKey(id string) (int, error) {
	if !strings.HasPrefix(id, "vm-") {
		return 0, fmt.Errorf("allocate: invalid vm id: %q", id)
	}
	n, err := strconv.Atoi(id[3:])
	if err != nil || n < 0 {
		return 0, fmt.Errorf("allocate: invalid vm id: %q", id)
	}
	return n, nil
}

func formatSlotKey(n int) string {
	return fmt.Sprintf("vm-%04d", n)
}
