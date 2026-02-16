package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/Amertz08/go-by-example/link"
	"github.com/Amertz08/go-by-example/link/kit/hlog"
	"github.com/Amertz08/go-by-example/link/kit/traceid"
	"github.com/Amertz08/go-by-example/link/rest"
)

type config struct {
	http struct {
		addr     string
		timeouts struct{ read, idle time.Duration }
	}
	lg *slog.Logger
}

func main() {
	var cfg config
	flag.StringVar(&cfg.http.addr, "http.addr", "localhost:8080", "http address to listen on")
	flag.DurationVar(
		&cfg.http.timeouts.read,
		"http.timeouts.read",
		20*time.Second,
		"http read timeout",
	)
	flag.DurationVar(
		&cfg.http.timeouts.idle,
		"http.timeouts.idle",
		40*time.Second,
		"http idle timeout",
	)
	flag.Parse()

	cfg.lg = slog.New(
		slog.NewTextHandler(os.Stderr, nil),
	).With("app", "linkd")
	cfg.lg.Info("starting", "addr", cfg.http.addr)

	if err := run(context.Background(), cfg); err != nil {
		cfg.lg.Error("failed to start server", "error", err)
		os.Exit(1)
	}
}

func run(_ context.Context, cfg config) error {
	shortener := new(link.Shortener)

	lg := slog.New(traceid.NewLogHandler(cfg.lg.Handler()))

	mux := http.NewServeMux()
	mux.Handle("POST /shorten", rest.Shorten(lg, shortener))
	mux.Handle("GET /{key}", rest.Resolve(lg, shortener))
	mux.HandleFunc("GET /health", rest.Health)

	loggingMiddleware := hlog.Middleware(lg)
	srv := &http.Server{
		Handler:     traceid.Middleware(loggingMiddleware(mux)),
		Addr:        cfg.http.addr,
		ReadTimeout: cfg.http.timeouts.read,
		IdleTimeout: cfg.http.timeouts.idle,
	}

	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("server closed unexpectedly: %w", err)
	}
	return nil
}
