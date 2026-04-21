package handler

import (
	"database/sql"
	"sync"
)

type Handler struct {
	db                  *sql.DB
	jobListeners        map[listener]struct{}
	jobListenersMu      sync.Mutex
	notebookListeners   map[listener]struct{}
	notebookListenersMu sync.Mutex
}

func New(db *sql.DB) *Handler {
	return &Handler{
		db:                db,
		jobListeners:      make(map[listener]struct{}),
		notebookListeners: make(map[listener]struct{}),
	}
}
