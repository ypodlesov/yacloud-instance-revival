package bg_process

import (
	"log/slog"
	"sync"
	"yacloud_revival/internal/config"
)

type BackgroundProcess struct {
	Config *config.Config
	Logger *slog.Logger
}

func (bg *BackgroundProcess) Run(wg *sync.WaitGroup) {
	bg.Logger.Info("bg_process started")

	var tasksWg sync.WaitGroup

	for _, instance := range bg.Config.Instances {
		t := &task{
			grpcAddress: bg.Config.Address,
			instanceId:  instance.InstanceId,
			period:      instance.CheckHealthPeriod,
			logger:      bg.Logger,
		}
		tasksWg.Add(1)
		go t.Run(&tasksWg)
	}

	tasksWg.Wait()

	bg.Logger.Info("bg_process finished")
	wg.Done()
}
