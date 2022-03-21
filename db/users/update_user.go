package users

import (
	"context"
	"github.com/swarnakumar/go-identity/db/sql/sqlc"
	"golang.org/x/crypto/bcrypt"
)

func (users Users) UpdateAttributes(ctx context.Context, email string, isAdmin, isActive bool, createdBy *string) (*sqlc.UpdateUserAttributesRow, error) {
	q := sqlc.New(users.db)

	params := sqlc.UpdateUserAttributesParams{
		IsActive: isActive,
		IsAdmin:  isAdmin,
		Email:    email,
	}

	newCtx := getCtxWithUser(ctx, createdBy)
	return q.UpdateUserAttributes(newCtx, params)

}

func (users Users) ChangePassword(ctx context.Context, email, password string, createdBy *string) (*sqlc.UpdateUserAttributesRow, error) {
	ok, complexityErr := CheckPassword(password)
	if !ok {
		return nil, complexityErr
	}

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return nil, SimplePasswordError
	}

	params := sqlc.UpdateUserAttributesParams{
		Passhash: string(hashedPwd),
		Email:    email,
	}

	q := sqlc.New(users.db)
	newCtx := getCtxWithUser(ctx, createdBy)
	return q.UpdateUserAttributes(newCtx, params)
}
