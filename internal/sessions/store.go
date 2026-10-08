package sessions

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type Message struct {
	Role    string    `json:"role"`
	Content string    `json:"content"`
	Time    time.Time `json:"time"`
}

type Session struct {
	ID         string    `json:"id"`
	Title      string    `json:"title"`
	Messages   []Message `json:"messages"`
	ProviderID string    `json:"provider_id"`
	Model      string    `json:"model"`
	Cost       float64   `json:"cost"`
	Tokens     int       `json:"tokens"`
}

type Store struct {
	sessions map[string]*Session
	storage  string
}

func NewStore() *Store {
	home, _ := os.UserHomeDir()
	path := filepath.Join(home, ".vefier_sessions")
	os.MkdirAll(path, 0755)
	
	return &Store{
		sessions: make(map[string]*Session),
		storage:  path,
	}
}

func (s *Store) CreateSession(provider, model string) *Session {
	id := time.Now().Format("20060102150405")
	sess := &Session{
		ID:         id,
		Title:      "Новый чат",
		ProviderID: provider,
		Model:      model,
		Messages:   []Message{},
	}
	s.sessions[id] = sess
	return sess
}

func (s *Store) AddMessage(id string, role, content string, tokens int, cost float64) {
	if sess, ok := s.sessions[id]; ok {
		sess.Messages = append(sess.Messages, Message{
			Role:    role,
			Content: content,
			Time:    time.Now(),
		})
		sess.Tokens += tokens
		sess.Cost += cost
		s.Save(id)
	}
}

func (s *Store) Save(id string) {
	sess, ok := s.sessions[id]
	if !ok {
		return
	}
	data, _ := json.Marshal(sess)
	os.WriteFile(filepath.Join(s.storage, id+".json"), data, 0644)
}

