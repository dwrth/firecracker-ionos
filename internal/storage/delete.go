package storage

import (
	"fmt"
	"os"
)

// DeleteSandboxStorage removes the sandbox storage directory for the given sandbox ID.
func (m *Manager) DeleteSandboxStorage(id string) error {
	dir := m.SandboxDir(id)

	if _, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("storage: stat sandbox dir: %w", err)
	}

	if err := os.RemoveAll(dir); err != nil {
		return fmt.Errorf("storage: delete sandbox storage: %w", err)
	}

	return nil
}
