package allocate_test

import (
	"testing"

	"github.com/dwrth/knaller/internal/allocate"
)

func TestInterfaceNames(t *testing.T) {
	got := allocate.InterfaceNames(4)
	want := allocate.Names{
		Namespace: "kn-sandbox-0004",
		HostVeth:  "kn4-host",
		NSVeth:    "kn4-ns",
		TAP:       "tap0",
	}
	if got != want {
		t.Fatalf("InterfaceNames() = %v, want %v", got, want)
	}
}
