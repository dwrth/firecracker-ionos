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

func TestStore_Create(t *testing.T) {
	dir := t.TempDir()
	s := state.New(dir)
	vm := state.VM{ID: "vm-0001", Name: "worker-1", Status: state.StatusStopped}

	if err := s.Create(vm); err != nil {
		t.Fatalf("Create() failed: %v", err)
	}
	if err := s.Create(vm); !errors.Is(err, state.ErrAlreadyExists) {
		t.Fatalf("Create() returned incorrect error: got %v, want %v", err, state.ErrAlreadyExists)
	}
}

func TestStore_List(t *testing.T) {
	dir := t.TempDir()
	s := state.New(dir)

	vms, err := s.List()
	if err != nil {
		t.Fatalf("List() failed: %v", err)
	}
	if len(vms) != 0 {
		t.Fatalf("List() returned incorrect number of VMs: got %d, want 0", len(vms))
	}

	if err := s.Create(state.VM{ID: "vm-0002", Name: "b"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Create(state.VM{ID: "vm-0001", Name: "a"}); err != nil {
		t.Fatal(err)
	}

	vms, err = s.List()
	if err != nil {
		t.Fatalf("List() failed: %v", err)
	}
	if len(vms) != 2 || vms[0].ID != "vm-0001" || vms[1].ID != "vm-0002" {
		t.Fatalf("List() returned incorrect VMs: got %+v, want [ID:vm-0001, Name:a, ID:vm-0002, Name:b]", vms)
	}
}

func TestStore_Delete(t *testing.T) {
	dir := t.TempDir()
	s := state.New(dir)
	if err := s.Create(state.VM{ID: "vm-0001"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Delete("vm-0001"); err != nil {
		t.Fatalf("Delete() failed: %v", err)
	}
	exists, err := s.Exists("vm-0001")
	if err != nil {
		t.Fatalf("Exists() failed: %v", err)
	}
	if exists {
		t.Fatalf("Exists() returned true for deleted VM")
	}
	if err := s.Delete("vm-0001"); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("Delete() returned incorrect error: got %v, want %v", err, os.ErrNotExist)
	}
}
