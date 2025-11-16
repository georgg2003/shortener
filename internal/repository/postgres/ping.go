package postgres

import "context"

func (r *repository) Ping(ctx context.Context) error {
	r.logger.WithContext(ctx).Debug("sent ping to database")
	return r.db.Ping(ctx)
}
