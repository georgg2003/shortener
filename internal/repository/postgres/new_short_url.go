package postgres

import (
	"context"
	"errors"
)

var errFailedToInsertNewShortURL = errors.New("failed to insert new short url")

func (r *repository) NewShortURL(ctx context.Context, url string, shortID string) error {
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
