package postgres

import (
	"context"

	"github.com/georgg2003/shortener/internal/models"
	"github.com/georgg2003/shortener/pkg/utils"
)

func (r *repository) GetStats(ctx context.Context) (models.Stats, error) {
	conn, err := r.db.Acquire(ctx)
	if err != nil {
		err = utils.ErrWrap(err, errFailedToAcquireConnection.Error())
		return models.Stats{}, err
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return models.Stats{}, utils.ErrWrap(err, "failed to begin tx")
	}
	var urlsCount int64
	err = tx.QueryRow(ctx, "SELECT count(*) FROM url_entity").Scan(&urlsCount)
	if err != nil {
		return models.Stats{}, utils.ErrWrap(err, "failetd to select url_entity count")
	}

	var usersCount int64
	err = tx.QueryRow(ctx, "SELECT count(*) FROM users").Scan(&usersCount)
	if err != nil {
		return models.Stats{}, utils.ErrWrap(err, "failetd to select users count")
	}

	return models.Stats{
		URLsCount:  urlsCount,
		UsersCount: usersCount,
	}, nil
}

func (r *repository) GetUsersStats(ctx context.Context) (int, error) {
	conn, err := r.db.Acquire(ctx)
	if err != nil {
		err = utils.ErrWrap(err, errFailedToAcquireConnection.Error())
		return -1, err
	}
	defer conn.Release()

	var count int
	conn.QueryRow(ctx, "SELECT count(*) FROM users").Scan(&count)
	return count, nil
}
