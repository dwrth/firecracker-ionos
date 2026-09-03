package allocate_test

import (
	"testing"

	"github.com/dwrth/knaller/internal/allocate"
	"github.com/dwrth/knaller/internal/state"
	"github.com/oklog/ulid/v2"
)

func TestCreds(t *testing.T) {
	uid, gid := allocate.Creds(4, 12000, 12000)
	if uid != 12004 || gid != 12004 {
		t.Errorf("Creds() = %d, %d, want 12004, 12004", uid, gid)
	}
}

func TestNewSandboxID(t *testing.T) {
	uniqueID := func() string { return ulid.Make().String() }
	existing := []state.Sandbox{
		{Name: "test1", ID: uniqueID()},
		{Name: "test2", ID: uniqueID()},
		{Name: "test3", ID: uniqueID()},
	}

	newID, err := allocate.NewSandboxID(existing)
	if err != nil {
		t.Fatal(err)
	}
	_, err = ulid.Parse(newID)
	if err != nil {
		t.Errorf("NewSandboxID() returned invalid id: %s", newID)
	}

	existing = append(existing, state.Sandbox{Name: "test4", ID: newID})
	secondID, err := allocate.NewSandboxID(existing)
	if err != nil {
		t.Fatal(err)
	}
	_, err = ulid.Parse(secondID)
	if err != nil {
		t.Errorf("NewSandboxID() returned invalid id: %s", secondID)
	}
	if secondID == newID {
		t.Errorf("NewSandboxID() returned duplicate id: %s", secondID)
	}
}
