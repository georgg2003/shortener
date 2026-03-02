package postgres

import (
	"context"

	"github.com/georgg2003/shortener/internal/models"
	"github.com/georgg2003/shortener/pkg/utils"
)

func (r *repository) GetShortIDsBatch(ctx context.Context, entities []*models.URLEntity) (map[string]string, error) {
	conn, err := r.db.Acquire(ctx)
	if err != nil {
		err = utils.ErrWrap(err, errFailedToAcquireConnection.Error())
		return nil, err
	}
	defer conn.Release()

	urls := make([]string, len(entities))
	for i, e := range entities {
		urls[i] = e.OriginalURL
	}

	rows, err := conn.Query(ctx, `
		SELECT original_url, short_id
		FROM url_entity
		WHERE original_url = ANY($1)
	`, urls)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	m := make(map[string]string)
	for rows.Next() {
		var orig, shortID string
		if err = rows.Scan(&orig, &shortID); err != nil {
			return nil, err
		}
		m[orig] = shortID
	}

	return m, err
}
