package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5/pgtype"
	storage "github.com/rammyblog/monitor-bee/internal/storage/sql"
)

type createMonitorRequest struct {
	Name               string         `json:"name"`
	Url                string         `json:"url"`
	Method             string         `json:"method"`
	IntervalSeconds    int16          `json:"interval_seconds"`
	TimeoutSeconds     int16          `json:"timeout_seconds"`
	Status             string         `json:"status"`
	Headers            map[string]any `json:"headers"`
	Body               string         `json:"body"`
	ExpectedStatusCode int            `json:"expected_status_code"`
}

type monitorResponse struct {
	ID                 int32          `json:"id"`
	UserID             int32          `json:"user_id"`
	Name               string         `json:"name"`
	Url                string         `json:"url"`
	Method             string         `json:"method"`
	IntervalSeconds    int32          `json:"interval_seconds"`
	TimeoutSeconds     int32          `json:"timeout_seconds"`
	Status             string         `json:"status"`
	Headers            map[string]any `json:"headers,omitempty"`
	Body               string         `json:"body,omitempty"`
	ExpectedStatusCode int            `json:"expected_status_code,omitempty"`
	CreatedAt          string         `json:"created_at"`
	UpdatedAt          string         `json:"updated_at"`
}

func toMonitorResponse(mon storage.Monitor) (monitorResponse, error) {
	resp := monitorResponse{
		ID:              mon.ID,
		UserID:          mon.UserID,
		Name:            mon.Name,
		Url:             mon.Url,
		Method:          mon.Method,
		IntervalSeconds: mon.IntervalSeconds,
		TimeoutSeconds:  mon.TimeoutSeconds,
		Status:          mon.Status,
		CreatedAt:       mon.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:       mon.UpdatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
	}

	if len(mon.Headers) > 0 {
		var headers map[string]any
		if err := json.Unmarshal(mon.Headers, &headers); err != nil {
			return resp, err
		}
		resp.Headers = headers
	}

	if mon.Body.Valid {
		resp.Body = mon.Body.String
	}

	if mon.ExpectedStatusCode.Valid {
		resp.ExpectedStatusCode = int(mon.ExpectedStatusCode.Int32)
	}

	return resp, nil
}

func (r createMonitorRequest) Valid() error {
	if r.Name == "" {
		return errors.New("name is required")
	}

	if r.Url == "" {
		return errors.New("url is required")
	}

	if r.Method == "" {
		return errors.New("method is required")
	}

	if r.Status == "" {
		return errors.New("status is required")
	}

	if r.IntervalSeconds < 30 {
		return errors.New("interval_seconds must be at least 30")
	}

	if r.TimeoutSeconds < 5 {
		return errors.New("timeout_seconds must be at least 5")
	}

	if r.TimeoutSeconds >= r.IntervalSeconds {
		return errors.New("timeout_seconds must be less than interval_seconds")
	}

	return nil

}

// Create monitor
func (s *Server) handleCreateMonitor() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value("userID").(int)

		ctx := r.Context()
		_, err := s.store.GetUserByID(ctx, int32(userID))
		if err != nil {
			respondError(w, r, http.StatusNotFound, ErrUserNotFound)
			return
		}

		req, err := decodeValid[createMonitorRequest](r)

		if err != nil {
			respondError(w, r, http.StatusBadRequest, err)
			return
		}

		// Convert headers map to JSON bytes
		var headersJSON []byte
		if req.Headers != nil {
			headersJSON, err = json.Marshal(req.Headers)
			if err != nil {
				respondError(w, r, http.StatusBadRequest, errors.New("invalid headers format"))
				return
			}
		}

		mon, err := s.store.CreateMonitor(ctx, storage.CreateMonitorParams{
			UserID:             int32(userID),
			Name:               req.Name,
			Url:                req.Url,
			Method:             req.Method,
			IntervalSeconds:    int32(req.IntervalSeconds),
			TimeoutSeconds:     int32(req.TimeoutSeconds),
			ExpectedStatusCode: pgtype.Int4{Int32: int32(req.ExpectedStatusCode), Valid: true},
			Status:             req.Status,
			Headers:            headersJSON,
			Body:               pgtype.Text{String: req.Body, Valid: true},
		})

		if err != nil {
			respondError(w, r, http.StatusBadRequest, err)
			return
		}

		// Convert to response format
		resp, err := toMonitorResponse(mon)
		if err != nil {
			respondError(w, r, http.StatusInternalServerError, err)
			return
		}

		respondJSON(w, r, resp)
	})

}

// GetMonitorByID
func (s *Server) handleGetMonitorByID() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			respondError(w, r, http.StatusBadRequest, ErrInvalidId)
			return
		}

		userID := r.Context().Value("userID").(int)

		ctx := r.Context()

		mon, err := s.store.GetMonitorByID(ctx, storage.GetMonitorByIDParams{
			ID:     int32(id),
			UserID: int32(userID),
		})

		if err != nil {
			if isNotFound(err) {
				respondError(w, r, http.StatusNotFound, ErrNotFound)
				return
			}
			respondError(w, r, http.StatusBadRequest, err)
			return
		}

		// Convert to response format
		resp, err := toMonitorResponse(mon)
		if err != nil {
			respondError(w, r, http.StatusInternalServerError, err)
			return
		}

		respondJSON(w, r, resp)

	})
}

//  ListMonitors,

// ListMonitorsByUser,

// ListActiveMonitors,
// ListMonitorsByStatus,

// ListMonitorsByUserAndStatus

// UpdateMonitor,
//
// UpdateMonitorStatus
// DeleteMonitor,
// DeleteMonitorByID
// MonitorExists,
//  UserOwnsMonitor,
//  CountMonitorsByUser,
//  CountActiveMonitorsByUser
