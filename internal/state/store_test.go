package state_test

import (
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
