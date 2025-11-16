package repository

import (
	"context"
	"errors"

	"github.com/georgg2003/shortener/internal/models"
)

var errFailedToInsertNewShortURL = errors.New("failed to insert new short url")

func (r *repository) newShortURLInStorage(ctx context.Context, url string, shortID string) error {
	r.storage.Store(shortID, models.ShortURL{
		ShortURL: shortID,
		LongURL:  url,
	})
	return nil
}

func (r *repository) newShortURLInDB(ctx context.Context, url string, shortID string) error {
	conn, err := r.db.Acquire(ctx)
	if err != nil {
		err = errors.Join(err, errFailedToAcquireConnection)
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, "INSERT INTO url_entity (short_id, original_url) VALUES ($1, $2)", shortID, url)
	if err != nil {
		err = errors.Join(err, errFailedToInsertNewShortURL)
		return err
	}

	return nil
}

func (r *repository) NewShortURL(ctx context.Context, url string, shortID string) error {
	if r.cfg.DataBaseDSN != "" {
		return r.newShortURLInDB(ctx, url, shortID)
	}
	return r.newShortURLInStorage(ctx, url, shortID)
}
