package repository

import "context"

func (r *repository) Ping(ctx context.Context) error {
	return r.db.Ping(ctx)
}
