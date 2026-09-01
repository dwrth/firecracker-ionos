package storage

import (
	"fmt"
	"os"
)

// DeleteVMStorage removes the VM storage directory for id.
func (m *Manager) DeleteVMStorage(id string) error {
	dir := m.VMDir(id)

	if _, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("storage: stat vm dir: %w", err)
	}

	if err := os.RemoveAll(dir); err != nil {
		return fmt.Errorf("storage: deletevm storage: %w", err)
	}

	return nil
}
