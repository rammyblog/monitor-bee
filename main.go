package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rammyblog/monitor-bee/internal/config"
	"github.com/rammyblog/monitor-bee/internal/handlers"
	"github.com/rammyblog/monitor-bee/internal/middleware"
	storage "github.com/rammyblog/monitor-bee/internal/storage/sql"
)

func main() {
	cfg := config.Load()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	// Initialize storage
	store, err := storage.NewStore(cfg.DatabaseURL)
	if err != nil {
		slog.Error("failed to initialize database", "error", err)
		os.Exit(1)
	}
	defer store.Close()

	// Setup handlers
	h := &handlers.Handler{
		Store:  store,
		Logger: logger,
	}

	// Setup router with middleware
	mux := http.NewServeMux()

	// Public routes
	mux.HandleFunc("GET /health", h.Health)
	mux.HandleFunc("POST /auth/login", h.Login)
	mux.HandleFunc("POST /auth/register", h.Register)

	// Protected routes
	protected := http.NewServeMux()
	protected.HandleFunc("GET /api/profile", h.GetProfile)
	protected.HandleFunc("PUT /api/profile", h.UpdateProfile)
	protected.HandleFunc("GET /api/users", h.ListUsers)

	// Apply auth middleware to protected routes
	mux.Handle("/api/", middleware.Auth(store)(protected))

	// Apply global middleware
	handler := middleware.CORS(
		middleware.Logging(logger)(
			middleware.Recovery(logger)(mux),
		),
	)

	// Setup server
	server := &http.Server{
		Addr:         cfg.Port,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		slog.Info("starting server", "port", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server failed to start", "error", err)
			os.Exit(1)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("server forced to shutdown", "error", err)
	}

	slog.Info("server exited")
}
