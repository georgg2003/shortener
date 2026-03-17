package postgres

import (
	"context"

	"github.com/georgg2003/shortener/internal/models"
	"github.com/georgg2003/shortener/pkg/utils"
)

func (r *repository) GetStats(ctx context.Context) (models.Stats, error) {
	var urlsCount, usersCount int64
	err := r.db.QueryRow(
		ctx,
		"SELECT (SELECT count(*) FROM url_entity), (SELECT count(*) FROM users)",
	).Scan(&urlsCount, &usersCount)
	if err != nil {
		return models.Stats{}, utils.ErrWrap(err, "failetd to select stats")
	}

	return models.Stats{
		URLsCount:  urlsCount,
		UsersCount: usersCount,
	}, nil
}
