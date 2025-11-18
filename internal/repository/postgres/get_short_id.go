package postgres

import (
	"context"
	"errors"

	repo "github.com/georgg2003/shortener/internal/repository"
	"github.com/georgg2003/shortener/pkg/utils"
	"github.com/jackc/pgx/v5"
)

func (r *repository) GetShortID(ctx context.Context, originalURL string) (string, error) {
	conn, err := r.db.Acquire(ctx)
	if err != nil {
		err = utils.ErrWrap(err, errFailedToAcquireConnection)
		return "", err
	}
	defer conn.Release()

	row := conn.QueryRow(ctx, "SELECT short_id FROM url_entity WHERE original_url = $1", originalURL)

	var shortID string
	err = row.Scan(&shortID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = utils.ErrWrap(err, repo.ErrNotFound)
		} else {
			err = utils.ErrWrap(err, errFailedToScan)
		}
		return "", err
	}

	return shortID, nil
}
