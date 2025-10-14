package server

import (
	"context"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(wrapped, r)

		s.logger.Info("request completed",
			"method", r.Method,
			"path", r.URL.Path,
			"status", wrapped.statusCode,
			"duration", time.Since(start),
			"remote_addr", r.RemoteAddr,
		)
	})
}

func (s *Server) recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				s.logger.Error("panic recovered",
					"error", err,
					"stack", string(debug.Stack()),
				)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *Server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			respondError(w, r, http.StatusUnauthorized, ErrUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			respondError(w, r, http.StatusUnauthorized, ErrInvalidToken)
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(s.jwtSecret), nil
		})

		if err != nil || !token.Valid {
			respondError(w, r, http.StatusUnauthorized, ErrInvalidToken)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			respondError(w, r, http.StatusUnauthorized, ErrInvalidToken)
			return
		}

		userID := int32(claims["user_id"].(float64))

		ctx := r.Context()
		exists, err := s.store.UserExists(ctx, userID)
		if err != nil || !exists {
			respondError(w, r, http.StatusUnauthorized, ErrUserNotFound)
			return
		}

		ctx = context.WithValue(ctx, "userID", int(userID))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
