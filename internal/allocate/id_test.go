package allocate_test

import (
	"testing"

	"github.com/dwrth/knaller/internal/allocate"
	"github.com/dwrth/knaller/internal/state"
	"github.com/oklog/ulid/v2"
)

func TestNextSlotKey(t *testing.T) {
	tests := []struct {
		name     string
		existing []state.Sandbox
		want     string
	}{
		{"empty", nil, "vm-0001"},
		{"sequential", []state.Sandbox{{ID: "vm-0001"}, {ID: "vm-0002"}}, "vm-0003"},
		{"hole", []state.Sandbox{{ID: "vm-0001"}, {ID: "vm-0003"}}, "vm-0002"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := allocate.NextSlotKey(tt.existing)
			if err != nil {
				t.Fatal(err)
			}
			if got != tt.want {
				t.Errorf("NextSlotKey() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestCreds(t *testing.T) {
	uid, gid, err := allocate.Creds("vm-0004", 12000, 12000)
	if err != nil {
		t.Fatal(err)
	}
	if uid != 12004 || gid != 12004 {
		t.Errorf("Creds() = %d, %d, want 12004, 12004", uid, gid)
	}

	_, _, err = allocate.Creds("bad-vm-id", 12000, 12000)
	if err == nil {
		t.Fatal("expected error for invalid id")
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
