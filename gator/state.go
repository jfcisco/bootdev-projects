package main

import (
	"database/sql"

	"github.com/jfcisco/gator/internal/config"
	"github.com/jfcisco/gator/internal/database"
)

type state struct {
	db     *database.Queries
	config *config.Config
}

func (s *state) LoadDb(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DbURL)
	if err != nil {
		return nil, err
	}
	s.db = database.New(db)
	return db, nil
}
