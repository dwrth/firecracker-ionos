package state_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/dwrth/knaller/internal/state"
)

func TestStore_Exists(t *testing.T) {
	dir := t.TempDir()
	s := state.New(dir)

	exists, err := s.Exists("vm-0001")
	if err != nil {
		t.Fatalf("Exists() failed: %v", err)
	}
	if exists {
		t.Fatalf("Exists() returned true for non-existent VM")
	}

	if err := os.WriteFile(filepath.Join(dir, "vm-0001.json"), []byte(`{}`), 0o644); err != nil {
		t.Fatal(err)
	}

	exists, err = s.Exists("vm-0001")
	if err != nil {
		t.Fatalf("Exists() failed: %v", err)
	}
	if !exists {
		t.Fatalf("Exists() returned false for existing VM")
	}
}

func TestStore_Save(t *testing.T) {
	dir := t.TempDir()
	s := state.New(dir)

	vm := state.VM{ID: "vm-0001", Name: "worker-1", Status: state.StatusStopped}
	if err := s.Save(vm); err != nil {
		t.Fatalf("Save() failed: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(dir, "vm-0001.json"))
	if err != nil {
		t.Fatal(err)
	}

	var got state.VM
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if got.ID != vm.ID || got.Name != vm.Name || got.Status != vm.Status {
		t.Fatalf("Save() returned incorrect VM: got %+v, want %+v", got, vm)
	}
}
