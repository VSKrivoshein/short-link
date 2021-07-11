package store

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Store struct {
	config         *Config
	db             *sql.DB
	userRepository *UserRepository
}

func NewStore(config *Config) *Store {
	return &Store{config: config}
}

func (s *Store) Open() error {
	db, err := sql.Open("postgres", s.config.DatabaseUrl)
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	s.db = db

	return nil
}

func (s *Store) Close() error {
	err := s.db.Close()
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) User() *UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
	}

	return s.userRepository
}
