package users

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"

	"github.com/swarnakumar/go-identity/db/sql/sqlc"
)

var UserAlreadyExistsError = errors.New("user already exists")

// Create a new user in the DB.
// Along the way also checks for:
//		1. User already existing.
//		2. Password complexity - raises if its too simple.
func (users Users) Create(
	ctx context.Context,
	email, password string,
	isAdmin bool,
	createdBy *string,
) (*sqlc.CreateUserRow, error) {

	newCtx := getCtxWithUser(ctx, createdBy)

	q := sqlc.New(users.db)
	if count, _ := q.CheckUserExists(newCtx, email); count == 1 {
		return nil, UserAlreadyExistsError
	}

	ok, complexityErr := CheckPassword(password)
	if !ok {
		return nil, complexityErr
	}

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return nil, SimplePasswordError
	}

	creds := sqlc.CreateUserParams{Email: email, Passhash: string(hashedPwd), IsAdmin: isAdmin, IsActive: true}
	user, err := q.CreateUser(newCtx, creds)
	if err != nil {
		return nil, err
	}

	return user, err

}
