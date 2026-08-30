package state_test

import (
	"errors"
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

func TestStore_Load(t *testing.T) {
	dir := t.TempDir()
	s := state.New(dir)

	vm := state.VM{ID: "vm-0001", Name: "worker-1", Status: state.StatusStopped}
	if err := s.Save(vm); err != nil {
		t.Fatalf("Save() failed: %v", err)
	}

	got, err := s.Load("vm-0001")
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}
	if got != vm {
		t.Fatalf("Load() returned incorrect VM: got %+v, want %+v", got, vm)
	}

	_, err = s.Load("vm-9999")
	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("Load() returned incorrect error: got %v, want %v", err, os.ErrNotExist)
	}
}
