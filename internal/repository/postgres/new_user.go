package postgres

import (
	"context"
	"errors"

	"github.com/georgg2003/shortener/pkg/utils"
)

var errFailedToInsertNewUser = errors.New("failed to insert new user")

func (r *repository) NewUser(ctx context.Context) (int64, error) {
	conn, err := r.db.Acquire(ctx)
	if err != nil {
		err = utils.ErrWrap(err, errFailedToAcquireConnection.Error())
		return 0, err
	}
	defer conn.Release()

	var userID int64

	err = conn.QueryRow(
		ctx,
		"INSERT INTO users DEFAULT VALUES RETURNING id",
	).Scan(&userID)

	if err != nil {
		err = utils.ErrWrap(err, errFailedToInsertNewUser.Error())
		return 0, err
	}

	return userID, nil
}
