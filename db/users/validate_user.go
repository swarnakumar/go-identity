package users

import (
	"context"
	"golang.org/x/crypto/bcrypt"

	"github.com/swarnakumar/go-identity/db/sql/sqlc"
)

func (users Users) VerifyPassword(ctx context.Context, email, password string) bool {
	q := sqlc.New(users.db)
	pwd, err := q.GetPasswordForUser(ctx, email)
	if err != nil {
		return false
	}
	err = bcrypt.CompareHashAndPassword([]byte(pwd), []byte(password))
	params := sqlc.RegisterLoginAttemptParams{
		Email:   email,
		Success: err == nil,
	}

	_ = q.RegisterLoginAttempt(ctx, params)
	return err == nil
}
