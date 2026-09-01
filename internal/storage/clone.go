package storage

import (
	"fmt"
	"os"
	"os/exec"
)

func (m *Manager) CloneBaseRootfs(id string) error {
	src := m.BaseRootfs()
	dst := m.Rootfs(id)

	if _, err := os.Stat(src); err != nil {
		return fmt.Errorf("storage: base rootfs: %w", err)
	}

	if _, err := os.Stat(dst); err == nil {
		return fmt.Errorf("storage: rootfs already exists: %s", dst)
	} else if !os.IsNotExist(err) {
		return err
	}

	if err := os.MkdirAll(m.VMDir(id), 0o755); err != nil {
		return fmt.Errorf("storage: mkdir vm dir: %w", err)
	}

	cmd := exec.Command("cp", "--reflink=auto", "--sparse=always", src, dst)
	if out, err := cmd.CombinedOutput(); err != nil {
		_ = os.Remove(dst)
		return fmt.Errorf("storage: clone rootfs: %w: %s", err, out)
	}

	return nil
}
