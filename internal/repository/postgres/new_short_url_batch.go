package postgres

import (
	"context"
	"errors"

	"github.com/georgg2003/shortener/internal/models"
	"github.com/georgg2003/shortener/pkg/postgres"
)

func (r *repository) NewShortURLBatch(ctx context.Context, entities []*models.URLEntity) error {
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

	var hasUniqueViolationErr bool

	for _, v := range entities {
		res, err := tx.Exec(ctx, `
			INSERT INTO url_entity (short_id, original_url)
			VALUES ($1, $2)
			ON CONFLICT DO NOTHING
		`, v.ShortID, v.OriginalURL)
		if err != nil {
			err = errors.Join(err, errFailedToInsertNewShortURL)
			return err
		}
		if res.RowsAffected() == 0 {
			hasUniqueViolationErr = true
		}
	}

	errTx := tx.Commit(ctx)
	if errTx != nil {
		return errTx
	}

	if hasUniqueViolationErr {
		return postgres.ErrUniqueViolation
	}

	return err
}
