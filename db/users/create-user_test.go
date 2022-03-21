package users

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/swarnakumar/go-identity/db/conn"
)

func TestCreateNewUser(t *testing.T) {
	ctx := context.Background()

	db := conn.NewPool(ctx)
	users := Users{db: db}

	email := "abc@xyz.com"
	pwd := "qwerty1234@96219"

	u, err := users.Create(ctx, email, pwd, false, nil)
	assert.Nil(t, err)
	assert.Equal(t, email, u.Email)
	assert.Equal(t, sql.NullString{
		String: "",
		Valid:  false,
	}, u.CreatedBy)
	assert.False(t, u.IsAdmin)
	assert.True(t, u.IsActive)

	// Try Again. Shouldn't get it - already exists
	u, err = users.Create(ctx, email, pwd, false, nil)
	assert.IsType(t, UserAlreadyExistsError, err)

	// Not with bad password
	email = "xyz@abc.com"
	pwd = "1"
	u, err = users.Create(ctx, email, pwd, false, nil)
	assert.IsType(t, SimplePasswordError, err)

	// Not with less bad password
	email = "xyz@abc.com"
	pwd = "1234567890"
	u, err = users.Create(ctx, email, pwd, false, nil)
	assert.IsType(t, SimplePasswordError, err)

	// Now with good password, and with created_by
	email = "xyz@abc.com"
	pwd = "wdcvoo@3454&*("
	createdBy := "swarnakumar"
	u, err = users.Create(ctx, email, pwd, false, &createdBy)
	assert.Nil(t, err)
	assert.Equal(t, createdBy, u.CreatedBy.String)
}
