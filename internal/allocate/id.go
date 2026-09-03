// Package allocate provides utilities generating unique sandbox identifiers.
package allocate

import (
	"fmt"

	"github.com/dwrth/knaller/internal/state"
	"github.com/oklog/ulid/v2"
)

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
