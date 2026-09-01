package storage

import "github.com/dwrth/knaller/internal/state"

// CustomizeMountedRootfsForTest writes guest files under mountPoint for tests.
func CustomizeMountedRootfsForTest(mountPoint string, vm state.VM) error {
	return customizeMountedRootfs(mountPoint, vm)
}
