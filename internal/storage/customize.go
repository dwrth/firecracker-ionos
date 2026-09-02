package storage

import (
	"crypto/rand"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/dwrth/knaller/internal/state"
)

// CustomizeRootfs writes per-sandbox identity and network settings into a cloned rootfs.
func (m *Manager) CustomizeRootfs(sandbox state.Sandbox) error {
	rootfs := m.Rootfs(sandbox.ID)
	if _, err := os.Stat(rootfs); err != nil {
		return fmt.Errorf("storage: rootfs: %w", err)
	}

	mountPoint, loop, err := mountRootfs(rootfs)
	if err != nil {
		return err
	}
	defer func() {
		_ = unmountRootfs(mountPoint, loop)
	}()

	if err := customizeMountedRootfs(mountPoint, sandbox); err != nil {
		return err
	}

	return nil
}

func customizeMountedRootfs(mountPoint string, sandbox state.Sandbox) error {
	etcDir := filepath.Join(mountPoint, "etc")
	if err := os.MkdirAll(etcDir, 0o755); err != nil {
		return fmt.Errorf("storage: mkdir etc: %w", err)
	}

	if err := os.WriteFile(
		filepath.Join(etcDir, "hostname"),
		[]byte(sandbox.ID+"\n"),
		0o644,
	); err != nil {
		return fmt.Errorf("storage: write hostname: %w", err)
	}

	host := fmt.Sprintf(
		"127.0.0.1 localhost\n127.0.1.1 %s\n\n::1 localhost ip6-localhost ip6-loopback\n",
		sandbox.ID,
	)
	if err := os.WriteFile(
		filepath.Join(etcDir, "hosts"),
		[]byte(host),
		0o644,
	); err != nil {
		return fmt.Errorf("storage: write hosts: %w", err)
	}

	networkDir := filepath.Join(etcDir, "systemd/network")
	if err := os.MkdirAll(networkDir, 0o755); err != nil {
		return fmt.Errorf("storage: mkdir network dir: %w", err)
	}

	network := fmt.Sprintf(`[Match]
Name=eth0
[Network]
Address=%s/30
Gateway=%s
IPv6AcceptRA=no
LinkLocalAddressing=no
`, sandbox.GuestIP, sandbox.GatewayIP)
	if err := os.WriteFile(
		filepath.Join(networkDir, "10-eth0.network"),
		[]byte(network),
		0o644,
	); err != nil {
		return fmt.Errorf("storage: write network config: %w", err)
	}

	machineID, err := newMachineID()
	if err != nil {
		return fmt.Errorf("storage: generate machine-id: %w", err)
	}
	if err := os.WriteFile(
		filepath.Join(etcDir, "machine-id"),
		[]byte(machineID+"\n"),
		0o644,
	); err != nil {
		return fmt.Errorf("storage: write machine-id: %w", err)
	}

	return nil
}

func newMachineID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}

func mountRootfs(imagePath string) (mountPoint, loop string, err error) {
	mountPoint, err = os.MkdirTemp("", "knaller-roofs-")
	if err != nil {
		return "", "", fmt.Errorf("storage: mkdir mount point: %w", err)
	}

	out, err := exec.Command("losetup", "--find", "--show", imagePath).Output()
	if err != nil {
		_ = os.RemoveAll(mountPoint)
		return "", "", fmt.Errorf("storage: losetup: %w", err)
	}
	loop = strings.TrimSpace(string(out))

	if out, err := exec.Command("mount", loop, mountPoint).CombinedOutput(); err != nil {
		_ = exec.Command("losetup", "-d", loop).Run()
		_ = os.RemoveAll(mountPoint)
		return "", "", fmt.Errorf("storage: mount: %w: %s", err, out)
	}

	return mountPoint, loop, nil
}

func unmountRootfs(mountPoint, loop string) error {
	if out, err := exec.Command("umount", mountPoint).CombinedOutput(); err != nil {
		return fmt.Errorf("storage: umount: %w: %s", err, out)
	}
	if out, err := exec.Command("losetup", "-d", loop).CombinedOutput(); err != nil {
		return fmt.Errorf("storage: losetup -d: %w: %s", err, out)
	}
	return os.RemoveAll(mountPoint)
}
