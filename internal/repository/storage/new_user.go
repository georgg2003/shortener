package storage

import (
	"context"

	repo "github.com/georgg2003/shortener/internal/repository"
)

func (r *repository) NewUser(ctx context.Context) (int64, error) {
	r.logger.Fatal("NewUser method is not implemented")
	return 0, repo.ErrNotImplemented
}
