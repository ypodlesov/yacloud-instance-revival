package bg_process

import (
	"context"
	"log/slog"
	"sync"
	"yacloud_revival/internal/config"
	"yacloud_revival/internal/token"
)

type BackgroundProcess struct {
	Config      *config.Config
	Logger      *slog.Logger
	TokenGetter *token.Token
}

func (bg *BackgroundProcess) Run(ctx context.Context, wg *sync.WaitGroup) {
	bg.Logger.Info("bg_process started")

	for _, instance := range bg.Config.Instances {
		t := &task{
			grpcAddress: bg.Config.Address,
			instanceId:  instance.InstanceId,
			period:      instance.CheckHealthPeriod,
			logger:      bg.Logger,
			tokenGetter: bg.TokenGetter,
		}
		wg.Add(1)
		go func() {
			t.Run(ctx, wg)
			wg.Done()
		}()
	}

	bg.Logger.Info("bg_process started tasks and finished")
}
