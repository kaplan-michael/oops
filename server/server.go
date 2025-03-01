package server

import (
	"context"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kaplan-michael/oops/errorpage"
	"go.uber.org/zap"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/kaplan-michael/oops/internal/config"
)

func Run(ctx context.Context, conf *config.Config) error {

	ep := errorpage.NewErrorPage(conf.Template, conf.Errors)

	// Set up the router.
	r := chi.NewRouter()

	// Set up middleware.
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(ZapLogger)
	r.Use(middleware.Recoverer)

	// Health check endpoint.
	r.Get("/healthz", healthzHandler)

	// Handle error pages. Since the ingress forwards errors with headers,
	// we use a catch-all endpoint here.
	r.Get("/*", errorHandler(ep))

	srv := &http.Server{
		Addr:    ":" + conf.Port,
		Handler: r,
	}

	// Run the server in a separate goroutine.
	go func() {
		zap.S().Infof("Starting error service on :%s", conf.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.S().Fatalf("Server failed: %v", err)
		}
	}()

	// Block until the provided context is cancelled.
	<-ctx.Done()
	zap.S().Info("Context cancelled, shutting down server...")

	// Shutdown the server using the provided context (which already has a timeout).
	return srv.Shutdown(ctx)

}

// ZapLogger is a custom middleware that logs HTTP requests using the global Zap logger.
func ZapLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		// Process the request.
		next.ServeHTTP(w, r)
		duration := time.Since(start)
		zap.S().Info("HTTP request",
			zap.String("method", r.Method),
			zap.String("url", r.URL.String()),
			zap.String("remote", r.RemoteAddr),
			zap.Duration("duration", duration),
		)
	})
}
