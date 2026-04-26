package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// PgxPool is the minimal subset of *pgxpool.Pool methods used by repositories.
// Both *pgxpool.Pool and pgxmock.PgxPoolIface satisfy this interface, which
// allows production code to use a real connection pool while tests inject a
// mock pool.
type PgxPool interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}
