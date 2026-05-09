package store

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"

	"steamshopemulator/internal/domain"
)

type Store struct {
	mu       sync.RWMutex
	dataPath string
	value    domain.Storefront
}

func New(dataPath string) (*Store, error) {
	s := &Store{dataPath: dataPath}
	if err := s.load(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Store) load() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	file, err := os.Open(s.dataPath)
	if err != nil {
		return err
	}
	defer file.Close()

	var current domain.Storefront
	if err := json.NewDecoder(file).Decode(&current); err != nil {
		return err
	}

	s.value = current
	return nil
}

func (s *Store) Get() domain.Storefront {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.value
}

func (s *Store) Save(next domain.Storefront) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := os.MkdirAll(filepath.Dir(s.dataPath), 0o755); err != nil {
		return err
	}

	payload, err := json.MarshalIndent(next, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(s.dataPath, payload, 0o644); err != nil {
		return err
	}

	s.value = next
	return nil
}

func EnsureDefault(dataPath string, payload domain.Storefront) error {
	if _, err := os.Stat(dataPath); err == nil {
		return nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(dataPath), 0o755); err != nil {
		return err
	}

	encoded, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(dataPath, encoded, 0o644)
}
