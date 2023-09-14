package secretaty

import (
	"Testovoe1/internal/http"
	"Testovoe1/internal/logic"
	"Testovoe1/internal/storage"
	"Testovoe1/internal/watcher"
	"context"
	"github.com/go-co-op/gocron"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

type Secretary struct {
	storage        *storage.Storage
	httpController *http.HttpHandlerController
	watcher        *watcher.Wathcer
	cron           *gocron.Scheduler
	logger         log.Logger
}

func NewSecretary(storage *storage.Storage, httpController *http.HttpHandlerController, watcher *watcher.Wathcer, cron *gocron.Scheduler, logger log.Logger) *Secretary {
	return &Secretary{storage: storage, httpController: httpController, watcher: watcher, cron: cron, logger: logger}
}

// Running a cron to collect new information, every 5 minutes
func (s *Secretary) StartCron() {
	s.cron.Every(5).Minute().Do(func() {
		s.logger.Log("job", "UpdateStatistics")
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		s.statisticsJob(ctx)
	})

	s.cron.StartAsync()
	s.cron.StartBlocking()
}

// The main function responsible for updating data
func (s *Secretary) statisticsJob(ctx context.Context) {
	comments, err := s.watcher.RequestPosts(ctx)
	if err != nil {
		level.Error(s.logger).Log("err", err)
	}
	statistics := logic.CreateStatistics(ctx, comments)
	if len(statistics) == 0 {
		return
	}
	for _, statistic := range statistics {
		func() {
			if err := s.storage.UpdateStatistics(ctx, statistic); err != nil {
				level.Error(s.logger).Log("err", err)
			}
		}()
	}
}

func (s *Secretary) StartHttp() {
	go s.httpController.StartHTTP()
}
