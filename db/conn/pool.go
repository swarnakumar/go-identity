package conn

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/swarnakumar/go-identity/config"
)

func NewPool(ctx context.Context) *pgxpool.Pool {

	pgConfig, err := pgxpool.ParseConfig(config.PgConnStr)

	if err != nil {
		panic(err)
	}

	pgConfig.BeforeAcquire = func(ctx context.Context, conn *pgx.Conn) bool {
		// set the tenant id into this connection's setting
		userId := ctx.Value("user_id")

		if userId != nil && userId != "" {
			_, err := conn.Exec(ctx, "select set_request_user($1)", userId)
			if err != nil {
				//panic(err) // or better to log the error, and then `return false` to destroy this connection instead of leaving it open.
				log.Printf(err.Error())
				return false
			}
		}
		return true
	}

	pgConfig.AfterRelease = func(conn *pgx.Conn) bool {
		// set the setting to be empty before this connection is released to pool
		_, err := conn.Exec(context.Background(), "select set_request_user($1)", "")
		if err != nil {
			log.Printf(err.Error())
			return false
		}
		return true
	}

	pool, err := pgxpool.ConnectConfig(ctx, pgConfig)

	return pool

}
