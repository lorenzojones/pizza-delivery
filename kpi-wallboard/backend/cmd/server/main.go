// Command server is the KPI wallboard backend: it refreshes Aircall and sales
// data on a schedule and serves a TV-ready dashboard JSON to the Svelte
// frontend.
package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"wallboard/internal/aircall"
	"wallboard/internal/api"
	"wallboard/internal/config"
	"wallboard/internal/refresh"
	"wallboard/internal/store"
)

// version is overridable at build time via -ldflags "-X main.version=...".
var version = "dev"

func main() {
	logger := log.New(os.Stdout, "[wallboard] ", log.LstdFlags|log.Lmsgprefix)

	cfg, err := config.Load()
	if err != nil {
		logger.Fatalf("config: %v", err)
	}
	logger.Printf("starting wallboard %s — title=%q tz=%s refresh=%s staff=%d demoAircall=%v",
		version, cfg.DashboardTitle, cfg.Timezone, cfg.RefreshInterval, len(cfg.Staff), cfg.DemoMode())

	st := store.New(cfg.StaleAfter)
	air := aircall.New(aircall.Options{
		APIID:    cfg.AircallAPIID,
		APIToken: cfg.AircallAPIToken,
		BaseURL:  cfg.AircallBaseURL,
		Logger:   logger,
	})
	rs := refresh.New(cfg, air, st, logger)

	// Root context cancelled on SIGINT/SIGTERM for graceful shutdown.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Run the refresh loop in the background (does an immediate first refresh).
	go rs.Run(ctx)

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           api.New(cfg, st, rs, logger, version).Handler(),
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		logger.Printf("http: listening on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatalf("http: %v", err)
		}
	}()

	<-ctx.Done()
	logger.Printf("shutting down…")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Printf("shutdown: %v", err)
	}
}
