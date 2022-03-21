package users

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/swarnakumar/go-identity/db/conn"
	"testing"
)

func TestVerifyPassword(t *testing.T) {
	ctx := context.Background()
	db := conn.NewPool(ctx)
	users := Users{db: db}

	email := "abc123@xyz.com"
	pwd := "qwerty1234fqmn@32496219"

	users.Create(ctx, email, pwd, false, nil)

	assert.True(t, users.VerifyPassword(ctx, email, pwd))

	// There should be ONE login attempt
	att, _ := users.GetLoginAttemptsForUser(ctx, email)
	assert.Equal(t, 1, len(att))
	assert.True(t, att[0].Success)
	assert.Equal(t, email, att[0].Email)

	// User's last_login time must also get updated
	u, _ := users.GetByEmail(ctx, email)
	assert.Equal(t, att[0].LoginTime, u.LastLogin)

	// Wrong Password
	pwd = "qwerty1234fqmn@32496219123"
	assert.False(t, users.VerifyPassword(ctx, email, pwd))
	att, _ = users.GetLoginAttemptsForUser(ctx, email)
	assert.Equal(t, 2, len(att))
	// LIFO
	assert.True(t, att[1].Success)
	assert.False(t, att[0].Success)

	// User's last_login time must NOT get updated
	u, _ = users.GetByEmail(ctx, email)
	assert.Equal(t, att[1].LoginTime, u.LastLogin)
	assert.NotEqual(t, att[0].LoginTime, u.LastLogin)

	// User that's not present
	email = "poiu@wert.com"
	assert.False(t, users.VerifyPassword(ctx, email, pwd))
}
