package postgres

import (
	"context"
	"database/sql"

	"github.com/georgg2003/shortener/internal/models"
	"github.com/georgg2003/shortener/pkg/contextlib"
	"github.com/georgg2003/shortener/pkg/postgres"
	"github.com/georgg2003/shortener/pkg/utils"
)

func (r *repository) NewShortURLBatch(ctx context.Context, entities []*models.URLEntity) error {
	conn, err := r.db.Acquire(ctx)
	if err != nil {
		err = utils.ErrWrap(err, errFailedToAcquireConnection.Error())
		return err
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		err = utils.ErrWrap(err, errFailedToBeginTransaction.Error())
		return err
	}

	var userIDParam sql.NullInt64
	userID, ok := contextlib.GetUserID(ctx)
	if ok && userID != 0 {
		userIDParam = sql.NullInt64{
			Int64: userID,
			Valid: true,
		}
	}

	var hasUniqueViolationErr bool

	for _, v := range entities {
		res, err := tx.Exec(ctx, `
			INSERT INTO url_entity (short_id, original_url, user_id)
			VALUES ($1, $2, $3)
			ON CONFLICT DO NOTHING
		`, v.ShortID, v.OriginalURL, userIDParam)
		if err != nil {
			err = utils.ErrWrap(err, errFailedToInsertNewShortURL.Error())
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
