package state

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// Store persists sandbox records as JSON files in a directory.
type Store struct {
	dir string
}

// ErrAlreadyExists is returned when Create is called for an existing sandbox ID.
var ErrAlreadyExists = errors.New("state: sandbox already exists")

// New returns a Store that reads and writes sandbox records under dir.
func New(dir string) *Store {
	return &Store{dir: dir}
}

func (s *Store) path(id string) string {
	return filepath.Join(s.dir, id+".json")
}

// Exists reports whether a sandbox record exists for the given sandbox ID.
func (s *Store) Exists(id string) (bool, error) {
	_, err := os.Stat(s.path(id))
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}

	return false, err
}

// Save writes sandbox to disk, replacing any existing record with the same sandbox ID.
func (s *Store) Save(sandbox Sandbox) error {
	now := time.Now().UTC()

	existing, err := s.Load(sandbox.ID)
	if err == nil {
		sandbox.CreatedAt = existing.CreatedAt
		sandbox.Generation = existing.Generation + 1
	} else if errors.Is(err, os.ErrNotExist) {
		sandbox.CreatedAt = now
		sandbox.Generation = 1
	} else {
		return err
	}

	sandbox.UpdatedAt = now

	if !sandbox.DesiredState.Valid() {
		return ErrInvalidState
	}
	if !sandbox.ObservedState.Valid() {
		return ErrInvalidState
	}

	if err := os.MkdirAll(s.dir, 0o755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(sandbox, "", "  ")
	if err != nil {
		return err
	}

	tmp, err := os.CreateTemp(s.dir, sandbox.ID+".*.tmp")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName)

	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Sync(); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}

	return os.Rename(tmpName, s.path(sandbox.ID))
}

// Load reads the sandbox record for the given sandbox ID.
func (s *Store) Load(id string) (Sandbox, error) {
	data, err := os.ReadFile(s.path(id))
	if err != nil {
		return Sandbox{}, err
	}

	var sandbox Sandbox
	if err := json.Unmarshal(data, &sandbox); err != nil {
		return Sandbox{}, err
	}

	return sandbox, nil
}

// Create writes a new sandbox record. It returns ErrAlreadyExists if the sandbox ID is taken.
func (s *Store) Create(sandbox Sandbox) error {
	exists, err := s.Exists(sandbox.ID)
	if err != nil {
		return err
	}
	if exists {
		return ErrAlreadyExists
	}
	return s.Save(sandbox)
}

// List returns all sandbox records sorted by sandbox ID.
func (s *Store) List() ([]Sandbox, error) {
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}

	var sandboxes []Sandbox
	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ".json" {
			continue
		}
		id := strings.TrimSuffix(e.Name(), ".json")
		sandbox, err := s.Load(id)
		if err != nil {
			return nil, err
		}
		sandboxes = append(sandboxes, sandbox)
	}
	sort.Slice(sandboxes, func(i, j int) bool { return sandboxes[i].ID < sandboxes[j].ID })
	return sandboxes, nil
}

// Delete removes the sandbox record for the given sandbox ID.
func (s *Store) Delete(id string) error {
	return os.Remove(s.path(id))
}
