package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/h-dav/vectr/internal/config"
	"github.com/h-dav/vectr/internal/database"
	"github.com/h-dav/vectr/internal/transport/rest"
)

const exitCodeFatal = 0

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Failed loading config: %v", err)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: config.AsSlogLevel(cfg.MinimumLogLevel),
	}))

	db, err := database.NewConnection(cfg, logger)
	if err != nil {
		logger.Error("Failed to create new database connection: %w", slog.Any("err", err))
		os.Exit(exitCodeFatal)
	}

	server, err := rest.NewServer(cfg, db, logger)
	if err != nil {
		logger.Error("failed making new rest server: %w", slog.Any("err", err))
		os.Exit(exitCodeFatal)
	}

	run(server, logger)
}

func run(server *http.Server, logger *slog.Logger) {
	wg := sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())

	wg.Add(1)

	go func(running *sync.WaitGroup) {
		defer running.Done()

		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

		select {
		case sig := <-signalChan:
			logger.Debug("received signal: %w", slog.Any("signal", sig))
			cancel()
		case <-ctx.Done():
			return
		}
	}(&wg)

	wg.Add(1)

	go func(running *sync.WaitGroup) {
		defer running.Done()

		logger.Info("starting rest server", "port", slog.String("Addr", server.Addr))

		err := server.ListenAndServe()
		if err != nil {
			return
		}
	}(&wg)

	<-ctx.Done()

	wg.Wait()
}
