package server

import (
	"net/http"
)

func (s *Server) handleHealth() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, r, map[string]string{"status": "ok"})
	})
}
