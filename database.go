package main

import (
	"database/sql"
	"sync"

	"github.com/ip812/blog/database"
	"github.com/ip812/blog/status"
)

type DBWrapper interface {
	Queries() (*database.Queries, error)
	DB() (*sql.DB, error)
}

type SwappableDB struct {
	mu    sync.RWMutex
	db    *sql.DB
	q     *database.Queries
	ready bool
}

func NewSwappableDB() *SwappableDB {
	return &SwappableDB{}
}

func (s *SwappableDB) Swap(db *sql.DB, queries *database.Queries) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.db = db
	s.q = queries
	s.ready = true
}

func (s *SwappableDB) Queries() (*database.Queries, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if !s.ready {
		return nil, status.ErrDatabaseNotReady
	}
	return s.q, nil
}

func (s *SwappableDB) DB() (*sql.DB, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if !s.ready {
		return nil, status.ErrDatabaseNotReady
	}
	return s.db, nil
}
