package repository

import (
	"context"
	"errors"

	"github.com/georgg2003/shortener/internal/models"
)

var errFailedToInsertNewShortURL = errors.New("failed to insert new short url")

func (r *repository) newShortURLInStorage(ctx context.Context, url string, shortURL string) error {
	r.storage.Store(shortURL, models.ShortURL{
		ShortURL: shortURL,
		LongURL:  url,
	})
	return nil
}

func (r *repository) newShortURLInDB(ctx context.Context, url string, shortURL string) error {
	conn, err := r.db.Acquire(ctx)
	if err != nil {
		err = errors.Join(err, errFailedToAcquireConnection)
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, "INSERT INTO short_url (short_url, long_url) VALUES ($1, $2)", shortURL, url)
	if err != nil {
		err = errors.Join(err, errFailedToInsertNewShortURL)
		return err
	}

	return nil
}

func (r *repository) NewShortURL(ctx context.Context, url string, shortURL string) error {
	if r.cfg.DataBaseDSN != "" {
		return r.newShortURLInDB(ctx, url, shortURL)
	}
	return r.newShortURLInStorage(ctx, url, shortURL)
}
