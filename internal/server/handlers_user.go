package server

import (
	"errors"
	"net/http"

	storage "github.com/rammyblog/monitor-bee/internal/storage/sql"
)

type updateProfileRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (r updateProfileRequest) Valid() error {
	if r.Name == "" {
		return errors.New("name is required")
	}
	if r.Email == "" {
		return errors.New("email is required")
	}
	return nil
}

func (s *Server) handleGetProfile() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("userID").(int)

		ctx := r.Context()
		user, err := s.store.GetUserByID(ctx, int32(userID))
		if err != nil {
			respondError(w, r, http.StatusNotFound, ErrUserNotFound)
			return
		}

		respondJSON(w, r, user)
	})
}

func (s *Server) handleUpdateProfile() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("userID").(int)

		req, err := decodeValid[updateProfileRequest](r)
		if err != nil {
			respondError(w, r, http.StatusBadRequest, err)
			return
		}

		ctx := r.Context()
		err = s.store.UpdateUser(ctx, storage.UpdateUserParams{
			ID:    int32(userID),
			Name:  req.Name,
			Email: req.Email,
		})
		if err != nil {
			s.logger.Error("failed to update profile", "error", err)
			respondError(w, r, http.StatusInternalServerError, errors.New("internal server error"))
			return
		}

		noContent(w, r)
	})
}

func (s *Server) handleListUsers() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		users, err := s.store.ListUsers(ctx)
		if err != nil {
			s.logger.Error("failed to list users", "error", err)
			respondError(w, r, http.StatusInternalServerError, errors.New("internal server error"))
			return
		}

		respondJSON(w, r, users)
	})
}
