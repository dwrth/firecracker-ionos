package allocate_test

import (
	"testing"

	"github.com/dwrth/knaller/internal/allocate"
)

func TestInterfaceNames(t *testing.T) {
	got, err := allocate.InterfaceNames("vm-0004")
	if err != nil {
		t.Fatal(err)
	}
	want := allocate.Names{
		Namespace: "kn-vm-0004",
		HostVeth:  "kn4-host",
		NSVeth:    "kn4-ns",
		TAP:       "tap0",
	}
	if got != want {
		t.Fatalf("InterfaceNames() = %v, want %v", got, want)
	}
}
