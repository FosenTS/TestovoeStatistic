package watcher

import (
	"Testovoe1/internal/entity"
	"Testovoe1/internal/storage"
	"context"
	"encoding/json"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"net/http"
)

type Wathcer struct {
	storage *storage.Storage
	logger  log.Logger
}

func NewWatcher(logger log.Logger) (*Wathcer, error) {
	return &Wathcer{
		logger: logger,
	}, nil
}

func (w *Wathcer) RequestPosts(ctx context.Context) ([]*entity.Comment, error) {
	log := log.With(w.logger, "method", "request Posts")
	api := "https://jsonplaceholder.typicode.com/comments"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, api, nil)
	if err != nil {
		level.Error(log).Log("err", err)
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		level.Error(log).Log("err", err)
		return nil, err
	}
	var entityComments []*entity.Comment
	if err := json.NewDecoder(res.Body).Decode(&entityComments); err != nil {
		level.Error(log).Log("err", err)
		return nil, err
	}
	return entityComments, nil
}
