package main

import (
	"database/sql"
	"sync"

	"github.com/ip812/blog/status"
)

type DBWrapper interface {
	DB() (*sql.DB, error)
}

type SwappableDB struct {
	mu    sync.RWMutex
	db    *sql.DB
	ready bool
}

func NewSwappableDB() *SwappableDB {
	return &SwappableDB{}
}

func (s *SwappableDB) Swap(db *sql.DB) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.db = db
	s.ready = true
}

func (s *SwappableDB) DB() (*sql.DB, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if !s.ready {
		return nil, status.ErrDatabaseNotReady
	}
	return s.db, nil
}
