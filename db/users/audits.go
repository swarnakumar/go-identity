package users

import (
	"context"
	"github.com/swarnakumar/go-identity/db/sql/sqlc"
)

func (users Users) GetLoginAttemptsForUser(ctx context.Context, email string, cursorArgs ...int32) ([]*sqlc.LoginAttempt, error) {

	var offset int32 = 0
	var limit int32 = 10

	if len(cursorArgs) > 0 {
		offset = cursorArgs[0]
		if len(cursorArgs) > 1 {
			limit = cursorArgs[1]
		}
	}

	params := sqlc.GetLoginAttemptsForUserParams{
		Email:  email,
		Limit:  limit,
		Offset: offset,
	}

	q := sqlc.New(users.db)
	return q.GetLoginAttemptsForUser(ctx, params)
}

func (users Users) GetChangesForUser(ctx context.Context, email string, cursorArgs ...int32) ([]*sqlc.UserChange, error) {
	var offset int32 = 0
	var limit int32 = 10

	if len(cursorArgs) > 0 {
		offset = cursorArgs[0]
		if len(cursorArgs) > 1 {
			limit = cursorArgs[1]
		}
	}

	params := sqlc.GetChangesForUserParams{
		Email:  email,
		Limit:  limit,
		Offset: offset,
	}

	q := sqlc.New(users.db)
	return q.GetChangesForUser(ctx, params)
}

func (users Users) GetChangesCount(ctx context.Context, email *string) (int64, error) {
	q := sqlc.New(users.db)

	if email == nil {
		return q.GetChangesCount(ctx)
	} else {
		return q.GetChangesCountForUser(ctx, *email)
	}
}

func (users Users) GetDeletions(ctx context.Context, cursorArgs ...int32) ([]*sqlc.UserDeletion, error) {
	var offset int32 = 0
	var limit int32 = 10

	if len(cursorArgs) > 0 {
		offset = cursorArgs[0]
		if len(cursorArgs) > 1 {
			limit = cursorArgs[1]
		}
	}

	params := sqlc.GetDeletionsParams{
		Limit:  limit,
		Offset: offset,
	}

	q := sqlc.New(users.db)
	return q.GetDeletions(ctx, params)
}

func (users Users) GetDeletionsCount(ctx context.Context) int64 {
	q := sqlc.New(users.db)

	count, err := q.GetDeletionsCount(ctx)
	if err == nil {
		return count
	} else {
		return 0
	}
}
