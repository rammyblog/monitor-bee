package server

import (
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	storage "github.com/rammyblog/monitor-bee/internal/storage/sql"
	"golang.org/x/crypto/bcrypt"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r loginRequest) Valid() error {
	if r.Email == "" {
		return errors.New("email is required")
	}
	if r.Password == "" {
		return errors.New("password is required")
	}
	return nil
}

type registerRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (r registerRequest) Valid() error {
	if r.Email == "" {
		return errors.New("email is required")
	}
	if r.Name == "" {
		return errors.New("name is required")
	}
	if r.Password == "" {
		return errors.New("password is required")
	}
	if len(r.Password) < 6 {
		return errors.New("password must be at least 6 characters")
	}
	return nil
}

type authResponse struct {
	Token string       `json:"token"`
	User  storage.User `json:"user"`
}

func (s *Server) handleLogin() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req, err := decodeValid[loginRequest](r)
		if err != nil {
			respondError(w, r, http.StatusBadRequest, err)
			return
		}

		ctx := r.Context()
		user, err := s.store.GetUser(ctx, req.Email)
		if err != nil {
			respondError(w, r, http.StatusUnauthorized, ErrInvalidCredentials)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
			respondError(w, r, http.StatusUnauthorized, ErrInvalidCredentials)
			return
		}

		token, err := s.generateToken(int(user.ID), user.Email)
		if err != nil {
			s.logger.Error("failed to generate token", "error", err)
			respondError(w, r, http.StatusInternalServerError, errors.New("internal server error"))
			return
		}

		respond(w, r, http.StatusOK, authResponse{
			Token: token,
			User:  user,
		})
	})
}

func (s *Server) handleRegister() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req, err := decodeValid[registerRequest](r)
		if err != nil {
			respondError(w, r, http.StatusBadRequest, err)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			s.logger.Error("failed to hash password", "error", err)
			respondError(w, r, http.StatusInternalServerError, errors.New("internal server error"))
			return
		}

		ctx := r.Context()
		user, err := s.store.CreateUser(ctx, storage.CreateUserParams{
			Email:    req.Email,
			Name:     req.Name,
			Password: string(hashedPassword),
		})
		if err != nil {
			s.logger.Error("failed to create user", "error", err)
			respondError(w, r, http.StatusConflict, errors.New("email already exists"))
			return
		}

		token, err := s.generateToken(int(user.ID), user.Email)
		if err != nil {
			s.logger.Error("failed to generate token", "error", err)
			respondError(w, r, http.StatusInternalServerError, errors.New("internal server error"))
			return
		}

		respond(w, r, http.StatusCreated, authResponse{
			Token: token,
			User:  user,
		})
	})
}

func (s *Server) generateToken(userID int, email string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}
