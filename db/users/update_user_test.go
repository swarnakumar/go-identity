package users

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/swarnakumar/go-identity/db/conn"
)

func TestUpdateUser(t *testing.T) {
	ctx := context.Background()

	db := conn.NewPool(ctx)
	users := Users{db: db}

	email := "update-user123@xyz.com"
	pwd := "qwerty1234fqmn@32496219"

	users.Create(ctx, email, pwd, false, nil)

	createdBy := "abc@xyz.com"
	_, err := users.UpdateAttributes(ctx, email, true, false, &createdBy)
	assert.Nil(t, err)

	changes, _ := users.GetChangesForUser(ctx, email)
	assert.Equal(t, 1, len(changes))

	// Now change again. Should NOT be any NEW changes!
	_, err = users.UpdateAttributes(ctx, email, true, false, &createdBy)

	changes, _ = users.GetChangesForUser(ctx, email)
	assert.Equal(t, 1, len(changes))

}

func TestChangePassword(t *testing.T) {
	ctx := context.Background()
	db := conn.NewPool(ctx)
	users := Users{db: db}

	email := "update-pwd345@xyz.com"
	pwd := "qwerty1234fqmn@32496219"

	users.Create(ctx, email, pwd, false, nil)

	newPwd := "qwerty1234fqmn@324219"
	_, err := users.ChangePassword(ctx, email, newPwd, nil)
	assert.Nil(t, err)

	changes, _ := users.GetChangesForUser(ctx, email)
	assert.Equal(t, 1, len(changes))

}
