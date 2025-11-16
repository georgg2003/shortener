package postgres

import (
	"context"
	"errors"

	repo "github.com/georgg2003/shortener/internal/repository"
	"github.com/jackc/pgx/v5"
)

var errWrongStorageType = errors.New("wrong type in storage")

func (r *repository) GetLongURL(ctx context.Context, shortID string) (string, error) {
	conn, err := r.db.Acquire(ctx)
	if err != nil {
		err = errors.Join(err, errFailedToAcquireConnection)
		return "", err
	}
	defer conn.Release()

	row := conn.QueryRow(ctx, "SELECT original_url FROM url_entity WHERE short_id = $1", shortID)

	var longURL string
	err = row.Scan(&longURL)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = errors.Join(err, repo.ErrShortURLNotFound)
		} else {
			err = errors.Join(err, errFailedToScan)
		}
		return "", err
	}

	return longURL, nil
}
