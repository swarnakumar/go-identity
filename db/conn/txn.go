package conn

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

func RunInTxn[T any](
	ctx context.Context,
	db *pgxpool.Pool,
	fn func() (T, error),
) (T, error) {
	var result T
	tx, err := db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return result, fmt.Errorf("unable to begin conn: %w", err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		} else {
			_ = tx.Commit(ctx)
		}
	}()

	res, fnErr := fn()
	err = fnErr
	return res, fnErr

}
