package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/georgg2003/shortener/pkg/contextlib"
)

var errFailedToInsertNewShortURL = errors.New("failed to insert new short url")

func (r *repository) NewShortURL(ctx context.Context, url string, shortID string) error {
	conn, err := r.db.Acquire(ctx)
	if err != nil {
		err = errors.Join(err, errFailedToAcquireConnection)
		return err
	}
	defer conn.Release()

	var userIDParam sql.NullInt64
	userID, ok := contextlib.GetUserID(ctx)
	if ok && userID != 0 {
		userIDParam = sql.NullInt64{
			Int64: userID,
			Valid: true,
		}
	}

	_, err = conn.Exec(ctx, "INSERT INTO url_entity (short_id, original_url, user_id) VALUES ($1, $2, $3)", shortID, url, userIDParam)
	if err != nil {
		err = errors.Join(err, errFailedToInsertNewShortURL)
		return err
	}

	return nil
}
