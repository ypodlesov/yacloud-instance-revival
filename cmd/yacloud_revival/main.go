package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"yacloud_revival/internal/background"
	"yacloud_revival/internal/config"
	"yacloud_revival/internal/token"
)

func main() {

	tokenGetter := &token.Token{}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	logger.Info("loading config ...")
	cfg := config.MustLoad()

	logger.Info(
		"config loaded, starting yacloud_revival service",
		slog.Any("config", *cfg),
	)

	bg := &bg_process.BackgroundProcess{
		Config:      cfg,
		Logger:      logger,
		TokenGetter: tokenGetter,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		bg.Run(ctx, &wg)
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan
	logger.Info("received signal, initiating graceful shutdown", slog.String("signal", sig.String()))

	cancel()
	logger.Info("called context cancel")

	shutdownTimeout := 15 * time.Second
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer shutdownCancel()

	tasksWait := make(chan struct{})
	go func() {
		wg.Wait()
		close(tasksWait)
	}()

	select {
	case <-tasksWait:
		logger.Info("all goroutines have finished")
	case <-shutdownCtx.Done():
		logger.Warn("shutdown timed out, exiting forcefully")
	}

	logger.Info("service has been stopped gracefully")
}
