package storage

import (
	"Testovoe1/internal/entity"
	"Testovoe1/pkg/db/postgresql"
	"context"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	log "github.com/sirupsen/logrus"
)

type Storage struct {
	pool postgresql.PGXPool
}

func (s *Storage) NewTicket(ctx context.Context) (pgx.Tx, error) {
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.RepeatableRead,
	})
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func NewStorage(ctx context.Context, PGXPool postgresql.PGXPool) (*Storage, error) {
	err := PGXPool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("couldn't ping: %w", err)
	}

	return &Storage{pool: PGXPool}, nil
}

// 2023-09-13 20:39:07.935803435 +0300
func (s *Storage) appendStatisticks(ctx context.Context, statistic *entity.Statistic) error {
	if _, err := s.pool.Exec(
		ctx,
		"INSERT INTO statisticks (postid, word, count, time) VALUES ($1, $2, $3, $4)",
		statistic.Postid,
		statistic.Word,
		statistic.Count,
		statistic.Time,
	); err != nil {
		log.Info(err)
		return err
	}

	return nil
}

func (s *Storage) UpdateStatisticks(ctx context.Context, statistic *entity.Statistic) error {
	var statisticSnaphot entity.Statistic
	if err := s.pool.QueryRow(
		ctx,
		"SELECT id, postId, word, count, time from statisticks where  postId = $1 and word = $2",
		statistic.Postid,
		statistic.Word,
	).Scan(
		&statisticSnaphot.Id,
		&statisticSnaphot.Postid,
		&statisticSnaphot.Word,
		&statisticSnaphot.Count,
		&statisticSnaphot.Time,
	); err != nil {
		if err == pgx.ErrNoRows {
			if err := s.appendStatisticks(ctx, statistic); err != nil {
				return err
			}
			return nil
		} else {
			return err
		}
	}
	if statisticSnaphot.Count == statistic.Count {
		return nil
	}
	if _, err := s.pool.Exec(
		ctx,
		"UPDATE statisticks SET count = $3, time = $4 WHERE postid = $1 and word = $2",
		statistic.Postid,
		statistic.Word,
		statistic.Count,
		statistic.Time,
	); err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetStatististicById(ctx context.Context, id int) ([]*entity.Statistic, error) {
	rows, err := s.pool.Query(
		ctx,
		"SELECT id, postid, word, count, time from statisticks where postid = $1",
		id,
	)
	if err != nil {
		return nil, err
	}
	var statistics []*entity.Statistic
	err = pgxscan.ScanAll(&statistics, rows)
	if err != nil {
		return nil, err
	}
	return statistics, err
}

func (s *Storage) GetAllStatistics(ctx context.Context) ([]*entity.Statistic, error) {
	rows, err := s.pool.Query(
		ctx,
		"SELECT id, postid, word, count, time from statisticks",
	)
	if err != nil {
		return nil, err
	}
	var statistics []*entity.Statistic
	err = pgxscan.ScanAll(&statistics, rows)
	if err != nil {
		return nil, err
	}
	return statistics, err
}
