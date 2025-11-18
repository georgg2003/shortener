package postgres

import (
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

var UniqueViolationError = errors.New("pg unique violation error")

func IsUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	isPgViolationErr := (errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation)
	return isPgViolationErr || errors.Is(err, UniqueViolationError)
}
