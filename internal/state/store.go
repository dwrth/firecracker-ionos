package state

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Store struct {
	dir string
}

var ErrAlreadyExists = errors.New("state: vm already exists")

func New(dir string) *Store {
	return &Store{dir: dir}
}

func (s *Store) path(id string) string {
	return filepath.Join(s.dir, id+".json")
}

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

func (s *Store) Save(vm VM) error {
	if err := os.MkdirAll(s.dir, 0o755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(vm, "", "  ")
	if err != nil {
		return err
	}

	tmp, err := os.CreateTemp(s.dir, vm.ID+".*.tmp")
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

	return os.Rename(tmpName, s.path(vm.ID))
}

func (s *Store) Load(id string) (VM, error) {
	data, err := os.ReadFile(s.path(id))
	if err != nil {
		return VM{}, err
	}

	var vm VM
	if err := json.Unmarshal(data, &vm); err != nil {
		return VM{}, err
	}

	return vm, nil
}

func (s *Store) Create(vm VM) error {
	exists, err := s.Exists(vm.ID)
	if err != nil {
		return err
	}
	if exists {
		return ErrAlreadyExists
	}
	return s.Save(vm)
}

func (s *Store) List() ([]VM, error) {
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}

	var vms []VM
	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ".json" {
			continue
		}
		id := strings.TrimSuffix(e.Name(), ".json")
		vm, err := s.Load(id)
		if err != nil {
			return nil, err
		}
		vms = append(vms, vm)
	}
	sort.Slice(vms, func(i, j int) bool { return vms[i].ID < vms[j].ID })
	return vms, nil
}

func (s *Store) Delete(id string) error {
	return os.Remove(s.path(id))
}
