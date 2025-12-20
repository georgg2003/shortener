package postgres

import (
	"context"
	"errors"

	repo "github.com/georgg2003/shortener/internal/repository"
	"github.com/georgg2003/shortener/pkg/utils"
	"github.com/jackc/pgx/v5"
)

func (r *repository) GetLongURL(ctx context.Context, shortID string) (string, bool, error) {
	conn, err := r.db.Acquire(ctx)
	if err != nil {
		err = utils.ErrWrap(err, errFailedToAcquireConnection.Error())
		return "", false, err
	}
	defer conn.Release()

	row := conn.QueryRow(ctx, "SELECT original_url, is_deleted FROM url_entity WHERE short_id = $1", shortID)

	var longURL string
	var isDeleted bool
	err = row.Scan(&longURL, &isDeleted)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = repo.ErrNotFound
		} else {
			err = utils.ErrWrap(err, errFailedToScan.Error())
		}
		return "", false, err
	}

	return longURL, isDeleted, nil
}
