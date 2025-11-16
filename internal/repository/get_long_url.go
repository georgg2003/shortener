package repository

import (
	"context"
	"errors"

	"github.com/georgg2003/shortener/internal/models"
	"github.com/jackc/pgx/v5"
)

var ErrShortURLNotFound = errors.New("value not found")
var errWrongStorageType = errors.New("wrong type in storage")

func (r *repository) getLongURLFromStorage(ctx context.Context, shortID string) (string, error) {
	storageValue, ok := r.storage.Load(shortID)
	if !ok {
		return "", ErrShortURLNotFound
	}

	shortURLModel, ok := storageValue.(models.ShortURL)
	if !ok {
		return "", errWrongStorageType
	}

	return shortURLModel.LongURL, nil
}

func (r *repository) getLongURLFromDB(ctx context.Context, shortURL string) (string, error) {
	conn, err := r.db.Acquire(ctx)
	if err != nil {
		err = errors.Join(err, errFailedToAcquireConnection)
		return "", err
	}
	defer conn.Release()

	row := conn.QueryRow(ctx, "SELECT long_url FROM short_url WHERE short_url = $1", shortURL)

	var longURL string
	err = row.Scan(&longURL)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = errors.Join(err, ErrShortURLNotFound)
		} else {
			err = errors.Join(err, errFailedToScan)
		}
		return "", err
	}

	return longURL, nil
}

func (r *repository) GetLongURL(ctx context.Context, shortID string) (string, error) {
	if r.cfg.PostgresEnabled {
		return r.getLongURLFromDB(ctx, shortID)
	}
	return r.getLongURLFromStorage(ctx, shortID)
}
