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
		t.Fatalf("Exists() returned true for non-existent sandbox")
	}

	if err := os.WriteFile(filepath.Join(dir, "vm-0001.json"), []byte(`{}`), 0o644); err != nil {
		t.Fatal(err)
	}

	exists, err = s.Exists("vm-0001")
	if err != nil {
		t.Fatalf("Exists() failed: %v", err)
	}
	if !exists {
		t.Fatalf("Exists() returned false for existing sandbox")
	}
}

func TestStore_Load(t *testing.T) {
	dir := t.TempDir()
	s := state.New(dir)

	sandbox := state.Sandbox{
		ID:            "vm-0001",
		Name:          "worker-1",
		DesiredState:  state.DesiredStopped,
		ObservedState: state.ObservedStopped,
	}
	if err := s.Save(sandbox); err != nil {
		t.Fatalf("Save() failed: %v", err)
	}

	got, err := s.Load("vm-0001")
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}
	if got.ID != sandbox.ID || got.Name != sandbox.Name || got.DesiredState != sandbox.DesiredState || got.ObservedState != sandbox.ObservedState {
		t.Fatalf("Load() returned incorrect sandbox: got %+v, want %+v", got, sandbox)
	}

	_, err = s.Load("vm-9999")
	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("Load() returned incorrect error: got %v, want %v", err, os.ErrNotExist)
	}
}

func TestStore_Create(t *testing.T) {
	dir := t.TempDir()
	s := state.New(dir)
	sandbox := state.Sandbox{
		ID:            "vm-0001",
		Name:          "worker-1",
		DesiredState:  state.DesiredStopped,
		ObservedState: state.ObservedStopped,
	}

	if err := s.Create(sandbox); err != nil {
		t.Fatalf("Create() failed: %v", err)
	}
	if err := s.Create(sandbox); !errors.Is(err, state.ErrAlreadyExists) {
		t.Fatalf("Create() returned incorrect error: got %v, want %v", err, state.ErrAlreadyExists)
	}

	got, err := s.Load("vm-0001")
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}
	if got.ID != sandbox.ID || got.Name != sandbox.Name || got.DesiredState != sandbox.DesiredState || got.ObservedState != sandbox.ObservedState {
		t.Fatalf("Create() created incorrect sandbox: got %+v, want %+v", got, sandbox)
	}
	if got.Generation != 1 {
		t.Fatalf("Create() set incorrect generation: got %d, want %d", got.Generation, sandbox.Generation+1)
	}

	sandboxWithMissingState := state.Sandbox{
		ID:            "vm-0002",
		Name:          "worker-1",
		ObservedState: state.ObservedStopped,
	}

	if err := s.Create(sandboxWithMissingState); err == nil {
		t.Fatalf("Create() did not return expected error: %v", state.ErrInvalidState)
	} else if !errors.Is(err, state.ErrInvalidState) {
		t.Fatalf("Create() did not return expected error: got: %v want: %v", err, state.ErrInvalidState)
	}

	sandboxWithInvalidState := state.Sandbox{
		ID:            "vm-0002",
		Name:          "worker-1",
		DesiredState:  "deleting",
		ObservedState: state.ObservedStopped,
	}

	if err := s.Create(sandboxWithInvalidState); err == nil {
		t.Fatalf("Create() did not return expected error: %v", state.ErrInvalidState)
	} else if !errors.Is(err, state.ErrInvalidState) {
		t.Fatalf("Create() did not return expected error: got: %v want: %v", err, state.ErrInvalidState)
	}
}

func TestStore_List(t *testing.T) {
	dir := t.TempDir()
	s := state.New(dir)

	sandboxes, err := s.List()
	if err != nil {
		t.Fatalf("List() failed: %v", err)
	}
	if len(sandboxes) != 0 {
		t.Fatalf("List() returned incorrect number of sandboxes: got %d, want 0", len(sandboxes))
	}

	if err := s.Create(state.Sandbox{ID: "vm-0002", Name: "b", DesiredState: state.DesiredStopped, ObservedState: state.ObservedStopped}); err != nil {
		t.Fatal(err)
	}
	if err := s.Create(state.Sandbox{ID: "vm-0001", Name: "a", DesiredState: state.DesiredStopped, ObservedState: state.ObservedStopped}); err != nil {
		t.Fatal(err)
	}

	sandboxes, err = s.List()
	if err != nil {
		t.Fatalf("List() failed: %v", err)
	}
	if len(sandboxes) != 2 || sandboxes[0].ID != "vm-0001" || sandboxes[1].ID != "vm-0002" {
		t.Fatalf("List() returned incorrect sandboxes: got %+v, want [ID:vm-0001, Name:a, ID:vm-0002, Name:b]", sandboxes)
	}
}

func TestStore_Delete(t *testing.T) {
	dir := t.TempDir()
	s := state.New(dir)
	if err := s.Create(state.Sandbox{ID: "vm-0001", DesiredState: state.DesiredStopped, ObservedState: state.ObservedStopped}); err != nil {
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
		t.Fatalf("Exists() returned true for deleted sandbox")
	}
	if err := s.Delete("vm-0001"); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("Delete() returned incorrect error: got %v, want %v", err, os.ErrNotExist)
	}
}

func TestStore_Save(t *testing.T) {
	dir := t.TempDir()
	s := state.New(dir)
	sandbox := state.Sandbox{
		ID:            "vm-0001",
		Name:          "worker-1",
		DesiredState:  state.DesiredStopped,
		ObservedState: state.ObservedStopped,
	}

	if err := s.Create(sandbox); err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	sandbox, err := s.Load("vm-0001")
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if err := s.Save(sandbox); err != nil {
		t.Fatalf("Save() failed: %v", err)
	}

	sandboxAfterSave, err := s.Load("vm-0001")
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if sandboxAfterSave.CreatedAt.Compare(sandbox.CreatedAt) != 0 {
		t.Fatalf("second Save() modified CreatedAt: got %v, want %v", sandboxAfterSave.CreatedAt, sandbox.CreatedAt)
	}
	if sandboxAfterSave.UpdatedAt.Compare(sandbox.UpdatedAt) != 1 {
		t.Fatalf("second Save() did not update UpdatedAt: %v", sandboxAfterSave.UpdatedAt)
	}

	if sandboxAfterSave.Generation != sandbox.Generation+1 {
		t.Fatalf("second Save() did not increment Generation: got %v, want %v", sandboxAfterSave.Generation, sandbox.Generation+1)
	}
}
