package server

import (
	"log/slog"
	"net/http"

	storage "github.com/rammyblog/monitor-bee/internal/storage/sql"
)

type Server struct {
	store     *storage.Store
	logger    *slog.Logger
	jwtSecret string
}

func NewServer(store *storage.Store, logger *slog.Logger, jwtSecret string) *Server {
	return &Server{
		store:     store,
		logger:    logger,
		jwtSecret: jwtSecret,
	}
}

func (s *Server) Handler() http.Handler {
	return s.routes()
}
