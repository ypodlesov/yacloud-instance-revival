package main

import (
	"log/slog"
	"os"
	"sync"
	bg_process "yacloud_revival/internal/background"
	"yacloud_revival/internal/config"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	logger.Info("loading config ...")
	cfg := config.MustLoad()

	logger.Info(
		"config loaded, starting yacloud_revival service",
		slog.Any("config", *cfg),
	)

	bg := &bg_process.BackgroundProcess{
		Config: cfg,
		Logger: logger,
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go bg.Run(&wg)

	wg.Wait()
}
