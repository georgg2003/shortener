package postgres

import (
	"context"

	"github.com/georgg2003/shortener/internal/models"
	"github.com/georgg2003/shortener/pkg/utils"
)

const UserURLsLimit = 100

func (r *repository) GetUserURLs(ctx context.Context, userID int64) ([]models.URLEntity, error) {
	conn, err := r.db.Acquire(ctx)
	if err != nil {
		err = utils.ErrWrap(err, errFailedToAcquireConnection.Error())
		return nil, err
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, `
		SELECT original_url, short_id
		FROM url_entity
		WHERE user_id = $1
		LIMIT $2
	`, userID, UserURLsLimit)
	if err != nil {
		return nil, utils.ErrWrap(err, "failed to get user urls")
	}

	arr := make([]models.URLEntity, 0, UserURLsLimit)
	for rows.Next() {
		var entity models.URLEntity
		if err = rows.Scan(&entity.OriginalURL, &entity.ShortID); err != nil {
			return nil, err
		}
		arr = append(arr, entity)
	}

	return arr, err
}
