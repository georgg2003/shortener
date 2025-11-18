package postgres

import (
	"context"

	"github.com/georgg2003/shortener/internal/models"
	repo "github.com/georgg2003/shortener/internal/repository"
	"github.com/georgg2003/shortener/pkg/utils"
)

func (r *repository) GetShortIDsBatch(ctx context.Context, entities []*models.URLEntity) error {
	conn, err := r.db.Acquire(ctx)
	if err != nil {
		err = utils.ErrWrap(err, errFailedToAcquireConnection)
		return err
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
		return err
	}
	defer rows.Close()

	m := make(map[string]string)
	for rows.Next() {
		var orig, shortID string
		if err := rows.Scan(&orig, &shortID); err != nil {
			return err
		}
		m[orig] = shortID
	}

	for _, e := range entities {
		if id, ok := m[e.OriginalURL]; ok {
			e.ShortID = id
		} else {
			return repo.ErrNotFound
		}
	}

	return nil
}
