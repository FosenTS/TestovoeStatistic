package main

import (
	"Testovoe1/internal/config"
	"Testovoe1/internal/http"
	"Testovoe1/internal/secretaty"
	"Testovoe1/internal/storage"
	"Testovoe1/internal/watcher"
	"Testovoe1/pkg/db/postgresql"
	"context"
	"github.com/go-co-op/gocron"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"os"
	"time"
)

func main() {
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = log.With(logger,
			"service", "statistics",
			"time:", log.DefaultTimestampUTC,
			"caller", log.DefaultCaller,
		)
	}

	config.LoadEnv()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	postgresCfg := config.Postgres()
	logger.Log("pgxpool", postgresCfg.WaitTimeout, postgresCfg.WaitTimeoutMilliseconds)
	pgxPool, err := postgresql.NewPGXPool(
		ctx,
		&postgresql.ClientConfig{
			MaxConnections:        50,
			MaxConnectionAttempts: 10,
			WaitingDuration:       postgresCfg.WaitTimeout,
			Username:              postgresCfg.Username,
			Password:              config.Env().PostgresPassword,
			Host:                  postgresCfg.Host,
			Port:                  postgresCfg.Port,
			DatabaseName:          postgresCfg.DatabaseName,
			UseCA:                 postgresCfg.UseCA,
			CaAbsPath:             postgresCfg.CAAbsPath,
			SimpleQueryMode:       postgresCfg.SimpleQueryMode,
		},
		logger,
	)
	if err != nil {
		level.Error(logger).Log("err", err)
	}

	storage, err := storage.NewStorage(ctx, pgxPool)
	if err != nil {
		level.Error(logger).Log("err", err)
	}

	httpHadl := http.NewHttpHandlerController(storage, logger)

	watcher, err := watcher.NewWatcher(logger)
	if err != nil {
		level.Error(logger).Log("err", err)
	}

	cron := gocron.NewScheduler(time.UTC)
	secr := secretaty.NewSecretary(storage, httpHadl, watcher, cron, logger)

	secr.StartHttp()
	secr.StartCron()
}
