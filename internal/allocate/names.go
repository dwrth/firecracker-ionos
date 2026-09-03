package allocate

import "fmt"

// Names holds Linux network interface names for a sandbox.
type Names struct {
	Namespace string
	HostVeth  string
	NSVeth    string
	TAP       string
}

// InterfaceNames returns network namespace and interface names for id.
func InterfaceNames(slot int) Names {
	return Names{
		Namespace: fmt.Sprintf("kn-sandbox-%04d", slot),
		HostVeth:  fmt.Sprintf("kn%d-host", slot),
		NSVeth:    fmt.Sprintf("kn%d-ns", slot),
		TAP:       "tap0",
	}
}
