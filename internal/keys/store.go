package keys

import "sync"

type Key struct {
	ID       string
	Name     string
	Value    string
	Provider string
	Status   string // active, invalid, limit
}

type Store struct {
	mu   sync.RWMutex
	Keys []Key
}

func NewStore() *Store {
	return &Store{
		Keys: []Key{
			{ID: "1", Name: "Personal OpenAI", Value: "sk-...", Provider: "openai", Status: "active"},
			{ID: "2", Name: "Backup OpenAI", Value: "sk-backup...", Provider: "openai", Status: "active"},
		},
	}
}

func (s *Store) Mu() *sync.RWMutex {
	return &s.mu
}

func (s *Store) AddKeys(newKeys []Key) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Keys = append(s.Keys, newKeys...)
}

func (s *Store) GetWorkingKey(providerID string) *Key {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for i, k := range s.Keys {
		if k.Provider == providerID && k.Status == "active" {
			return &s.Keys[i]
		}
	}
	return nil
}

func (s *Store) MarkFailed(id string, reason string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.Keys {
		if s.Keys[i].ID == id {
			s.Keys[i].Status = reason
			return
		}
	}
}

