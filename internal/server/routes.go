package server

import (
	"net/http"
)

func (s *Server) routes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /health", s.handleHealth())
	mux.Handle("POST /auth/login", s.handleLogin())
	mux.Handle("POST /auth/register", s.handleRegister())

	// Protected routes - apply auth middleware
	mux.Handle("GET /api/profile", s.authMiddleware(s.handleGetProfile()))
	mux.Handle("PUT /api/profile", s.authMiddleware(s.handleUpdateProfile()))
	mux.Handle("GET /api/users", s.authMiddleware(s.handleListUsers()))

	// Monitors
	mux.Handle("POST /api/monitors", s.authMiddleware(s.handleCreateMonitor()))
	mux.Handle("GET /api/monitors/{id}", s.authMiddleware(s.handleGetMonitorByID()))

	return s.corsMiddleware(
		s.loggingMiddleware(
			s.recoveryMiddleware(mux),
		),
	)
}
