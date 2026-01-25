package storage

import (
	"context"

	"github.com/georgg2003/shortener/internal/repository/db"
)

func (r *repository) NewUser(ctx context.Context) (int64, error) {
	r.logger.Fatal("NewUser method is not implemented")
	return 0, db.ErrNotImplemented
}
