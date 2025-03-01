package main

import (
	"context"
	"github.com/kaplan-michael/oops/internal/config"
	"github.com/kaplan-michael/oops/internal/logger"
	"github.com/kaplan-michael/oops/server"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	// Load the configuration.
	conf := config.MustGetConfig()

	// Initialize the logger.
	if err := logger.Init(conf.LogLevel); err != nil {
		panic(err)
	}
	// Ensure logs are flushed on exit.
	defer zap.L().Sync()
	defer zap.S().Sync()

	zap.S().Infof("Configuration loaded: TemplateFile=%s, ErrorsFile=%s, LogLevel=%s", conf.Template, conf.Errors, conf.LogLevel)
	zap.S().Info("Starting...")

	// Create a cancellable context.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up OS signal handling.
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Listen for termination signals in a separate goroutine.
	go func() {
		<-sigChan
		zap.S().Info("Received termination signal, shutting down...")
		cancel()
	}()

	// Call server run, which will block until the context is cancelled.
	if err := server.Run(ctx, conf); err != nil {
		zap.S().Fatalf("Server error: %v", err)
	}

	zap.S().Info("Server shutdown gracefully.")
}
