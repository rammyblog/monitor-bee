package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rammyblog/monitor-bee/internal/config"
	"github.com/rammyblog/monitor-bee/internal/server"
	storage "github.com/rammyblog/monitor-bee/internal/storage/sql"
)

func main() {
	if err := run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	cfg := config.Load()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	store, err := storage.NewStore(cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	defer store.Close()

	srv := server.NewServer(store, logger, cfg.JWTSecret)
	httpServer := &http.Server{
		Addr:         cfg.Port,
		Handler:      srv.Handler(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
		BaseContext:  func(_ net.Listener) context.Context { return context.Background() },
	}

	go func() {
		logger.Info("starting server", "port", cfg.Port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server failed to start", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	logger.Info("server exited")
	return nil
}
