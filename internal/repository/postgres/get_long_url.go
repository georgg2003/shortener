package postgres

import (
	"context"
	"errors"

	repo "github.com/georgg2003/shortener/internal/repository"
	"github.com/georgg2003/shortener/pkg/utils"
	"github.com/jackc/pgx/v5"
)

func (r *repository) GetLongURL(ctx context.Context, shortID string) (string, error) {
	conn, err := r.db.Acquire(ctx)
	if err != nil {
		err = utils.ErrWrap(err, errFailedToAcquireConnection.Error())
		return "", err
	}
	defer conn.Release()

	row := conn.QueryRow(ctx, "SELECT original_url FROM url_entity WHERE short_id = $1", shortID)

	var longURL string
	err = row.Scan(&longURL)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = repo.ErrNotFound
		} else {
			err = utils.ErrWrap(err, errFailedToScan.Error())
		}
		return "", err
	}

	return longURL, nil
}
