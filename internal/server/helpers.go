package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func decode[T any](r *http.Request) (T, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, fmt.Errorf("decode json: %w", err)
	}
	return v, nil
}

func encode[T any](w http.ResponseWriter, r *http.Request, status int, v T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}
	return nil
}

type validator interface {
	Valid() error
}

func decodeValid[T validator](r *http.Request) (T, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, fmt.Errorf("decode json: %w", err)
	}
	if err := v.Valid(); err != nil {
		return v, fmt.Errorf("validation error: %w", err)
	}
	return v, nil
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidToken       = errors.New("invalid token")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrInvalidId          = errors.New("invalid id")
	ErrNotFound           = errors.New("resource not found")
)

func respond(w http.ResponseWriter, r *http.Request, status int, data any) {
	if err := encode(w, r, status, data); err != nil {
		return
	}
}

func respondError(w http.ResponseWriter, r *http.Request, status int, err error) {
	respond(w, r, status, map[string]string{
		"error": err.Error(),
	})
}

func respondJSON(w http.ResponseWriter, r *http.Request, data any) {
	respond(w, r, http.StatusOK, data)
}

func noContent(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func limited(r io.Reader, n int64) ([]byte, error) {
	lr := io.LimitedReader{R: r, N: n + 1}
	data, err := io.ReadAll(&lr)
	if err != nil {
		return nil, err
	}
	if lr.N <= 0 {
		return nil, fmt.Errorf("request body too large")
	}
	return data, nil
}

// isNotFound checks if the error is a "no rows" error from the database
func isNotFound(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "no rows in result set")
}
