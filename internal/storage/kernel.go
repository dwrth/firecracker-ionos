package storage

import (
	"fmt"
	"io"
	"os"
)

func (m *Manager) InstallKernel(id string) error {
	src := m.KernelSource()
	dst := m.Kernel(id)

	if _, err := os.Stat(src); err != nil {
		return fmt.Errorf("storage: kernel source: %w", err)
	}

	if _, err := os.Stat(dst); err == nil {
		return fmt.Errorf("storage: kernel already exists: %s", dst)
	} else if !os.IsNotExist(err) {
		return err
	}

	if err := os.MkdirAll(m.VMDir(id), 0o755); err != nil {
		return fmt.Errorf("storage: mkdir vm dir: %w", err)
	}

	if err := copyFile(src, dst); err != nil {
		return fmt.Errorf("storage install kernel: %w", err)
	}

	return nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0o644)
	if err != nil {
		return err
	}

	if _, err := io.Copy(out, in); err != nil {
		_ = out.Close()
		_ = os.Remove(dst)
		return err
	}

	if err := out.Sync(); err != nil {
		_ = out.Close()
		_ = os.Remove(dst)
		return err
	}

	return out.Close()
}
