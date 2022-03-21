package users

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/swarnakumar/go-identity/db/sql/sqlc"
)

type Users struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Users {
	return &Users{db: db}
}

func (users Users) GetByEmail(ctx context.Context, email string) (*sqlc.GetUserByEmailRow, error) {
	q := sqlc.New(users.db)
	return q.GetUserByEmail(ctx, email)
}

func (users Users) Listing(ctx context.Context, offset, limit int32) ([]*sqlc.ListUsersRow, error) {
	params := sqlc.ListUsersParams{
		Limit:  limit,
		Offset: offset,
	}

	q := sqlc.New(users.db)
	return q.ListUsers(ctx, params)
}

func (users Users) CheckEmailExists(ctx context.Context, email string) bool {
	q := sqlc.New(users.db)
	c, err := q.CheckUserExists(ctx, email)
	if err == nil {
		return c == 1
	} else {
		return false
	}
}

func (users Users) Count(ctx context.Context) int64 {
	q := sqlc.New(users.db)
	count, err := q.GetUserCount(ctx)

	if err != nil {
		count = 0
	}
	return count
}

func (users Users) Delete(ctx context.Context, email string, deletedBy *string) error {
	q := sqlc.New(users.db)

	newCtx := getCtxWithUser(ctx, deletedBy)
	err := q.DeleteUser(newCtx, email)
	return err
}
