package storage

import (
	"context"
)

func (r *repository) GetShortID(ctx context.Context, originalURL string) (string, error) {
	r.logger.Fatal("GetShortID is not implemented")
	return "", nil
}
