package storage

import (
	"fmt"
	"os"
	"strings"
)

// ValidateSandboxStorage reports whether all expected files exist for id.
func (m *Manager) ValidateSandboxStorage(id string) error {
	required := []struct {
		name string
		path string
	}{
		{name: "rootfs", path: m.Rootfs(id)},
		{name: "kernel", path: m.Kernel(id)},
		{name: "config", path: m.Config(id)},
	}

	var missing []string
	for _, file := range required {
		if _, err := os.Stat(file.path); err != nil {
			if os.IsNotExist(err) {
				missing = append(missing, file.name)
				continue
			}
			return fmt.Errorf("storage: validate %s: %w", file.name, err)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("storage: missing sandbox storage for %s: %s", id, strings.Join(missing, ", "))
	}

	return nil
}
