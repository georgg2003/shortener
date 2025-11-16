package repository

import (
	"context"
	"errors"

	"github.com/georgg2003/shortener/internal/models"
)

// var errFailedToInsertNewShortURL = errors.New("failed to insert new short url")

func (r *repository) newShortURLBatchInStorage(ctx context.Context, entities []*models.URLEntity) error {
	for _, v := range entities {
		r.storage.Store(v.ShortID, models.ShortURL{
			ShortURL: v.ShortID,
			LongURL:  v.OriginalURL,
		})
	}
	return nil
}

func (r *repository) newShortURLBatchInDB(ctx context.Context, entities []*models.URLEntity) error {
	conn, err := r.db.Acquire(ctx)
	if err != nil {
		err = errors.Join(err, errFailedToAcquireConnection)
		return err
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		err = errors.Join(err, errFailedToBeginTransaction)
		return err
	}

	for _, v := range entities {
		_, err = tx.Exec(ctx, "INSERT INTO url_entity (short_id, original_url) VALUES ($1, $2)", v.ShortID, v.OriginalURL)
		if err != nil {
			err = errors.Join(err, errFailedToInsertNewShortURL)
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *repository) NewShortURLBatch(ctx context.Context, entities []*models.URLEntity) error {
	if r.cfg.DataBaseDSN != "" {
		return r.newShortURLBatchInDB(ctx, entities)
	}
	return r.newShortURLBatchInStorage(ctx, entities)
}
