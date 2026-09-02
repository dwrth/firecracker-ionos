package storage

import "github.com/dwrth/knaller/internal/state"

func CustomizeMountedRootfsForTest(mountPoint string, sandbox state.Sandbox) error {
	return customizeMountedRootfs(mountPoint, sandbox)
}
