package postgres

import (
	"context"

	"github.com/georgg2003/shortener/internal/models"
	"github.com/georgg2003/shortener/pkg/utils"
)

func (r *repository) DeleteUserURLs(ctx context.Context, tasks []models.DeleteUserURLsTask) error {
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

	for _, task := range tasks {
		tx.Exec(ctx, `
			UPDATE url_entity
			SET is_deleted = true
			WHERE short_id = ANY($1) AND user_id = $2
		`, task.ShortIDs, task.UserID)
		if err != nil {
			err = utils.ErrWrap(err, "failed to delete user urls")
			return err
		}
	}

	errTx := tx.Commit(ctx)
	if errTx != nil {
		return errTx
	}

	return nil
}
