package postgres

import (
	"context"

	"github.com/georgg2003/shortener/internal/models"
	"github.com/georgg2003/shortener/pkg/utils"
)

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
	`, userID)
	if err != nil {
		return nil, utils.ErrWrap(err, "failed to get user urls")
	}

	arr := make([]models.URLEntity, 0)
	for rows.Next() {
		var orig, shortID string
		if err := rows.Scan(&orig, &shortID); err != nil {
			return nil, err
		}
		arr = append(arr, models.URLEntity{
			OriginalURL: orig,
			ShortID:     shortID,
		})
	}

	return arr, err
}
