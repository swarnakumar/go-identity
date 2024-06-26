// Code generated by sqlc. DO NOT EDIT.
// source: users.sql

package sqlc

import (
	"context"
	"database/sql"
	"time"
)

const checkUserExists = `-- name: CheckUserExists :one
SELECT count(*)::int FROM users WHERE email = $1
`

func (q *Queries) CheckUserExists(ctx context.Context, email string) (int32, error) {
	row := q.db.QueryRow(ctx, checkUserExists, email)
	var column_1 int32
	err := row.Scan(&column_1)
	return column_1, err
}

const createUser = `-- name: CreateUser :one
with u AS (INSERT INTO users (email, passhash, is_active, is_admin, created_at, updated_at, created_by)
    VALUES ($1, $2, $3, $4,
            now(), now(), nullif(current_setting('request.user', true), '') )
    RETURNING email, created_by, passhash, is_active, is_admin, created_at, updated_at, last_login)
SELECT email, is_active, is_admin, created_at, updated_at, last_login, created_by from u
`

type CreateUserParams struct {
	Email    string `json:"email"`
	Passhash string `json:"passhash"`
	IsActive bool   `json:"is_active"`
	IsAdmin  bool   `json:"is_admin"`
}

type CreateUserRow struct {
	Email     string         `json:"email"`
	IsActive  bool           `json:"is_active"`
	IsAdmin   bool           `json:"is_admin"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	LastLogin sql.NullTime   `json:"last_login"`
	CreatedBy sql.NullString `json:"created_by"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (*CreateUserRow, error) {
	row := q.db.QueryRow(ctx, createUser,
		arg.Email,
		arg.Passhash,
		arg.IsActive,
		arg.IsAdmin,
	)
	var i CreateUserRow
	err := row.Scan(
		&i.Email,
		&i.IsActive,
		&i.IsAdmin,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.LastLogin,
		&i.CreatedBy,
	)
	return &i, err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM users WHERE email = $1
`

func (q *Queries) DeleteUser(ctx context.Context, email string) error {
	_, err := q.db.Exec(ctx, deleteUser, email)
	return err
}

const getPasswordForUser = `-- name: GetPasswordForUser :one
SELECT passhash from users where email = $1
`

func (q *Queries) GetPasswordForUser(ctx context.Context, email string) (string, error) {
	row := q.db.QueryRow(ctx, getPasswordForUser, email)
	var passhash string
	err := row.Scan(&passhash)
	return passhash, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT users.email, is_active, is_admin, created_at, updated_at, last_login, created_by
FROM users WHERE users.email = $1
`

type GetUserByEmailRow struct {
	Email     string         `json:"email"`
	IsActive  bool           `json:"is_active"`
	IsAdmin   bool           `json:"is_admin"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	LastLogin sql.NullTime   `json:"last_login"`
	CreatedBy sql.NullString `json:"created_by"`
}

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (*GetUserByEmailRow, error) {
	row := q.db.QueryRow(ctx, getUserByEmail, email)
	var i GetUserByEmailRow
	err := row.Scan(
		&i.Email,
		&i.IsActive,
		&i.IsAdmin,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.LastLogin,
		&i.CreatedBy,
	)
	return &i, err
}

const getUserCount = `-- name: GetUserCount :one
SELECT COUNT(*) FROM users
`

func (q *Queries) GetUserCount(ctx context.Context) (int64, error) {
	row := q.db.QueryRow(ctx, getUserCount)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const listUsers = `-- name: ListUsers :many
SELECT email, is_admin, is_active, created_at, updated_at, created_by, last_login
FROM users
ORDER BY COALESCE(NULLIF($3, created_at), created_at)
LIMIT $1 OFFSET $2
`

type ListUsersParams struct {
	Limit   int32       `json:"limit"`
	Offset  int32       `json:"offset"`
	Column3 interface{} `json:"column_3"`
}

type ListUsersRow struct {
	Email     string         `json:"email"`
	IsAdmin   bool           `json:"is_admin"`
	IsActive  bool           `json:"is_active"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	CreatedBy sql.NullString `json:"created_by"`
	LastLogin sql.NullTime   `json:"last_login"`
}

func (q *Queries) ListUsers(ctx context.Context, arg ListUsersParams) ([]*ListUsersRow, error) {
	rows, err := q.db.Query(ctx, listUsers, arg.Limit, arg.Offset, arg.Column3)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*ListUsersRow
	for rows.Next() {
		var i ListUsersRow
		if err := rows.Scan(
			&i.Email,
			&i.IsAdmin,
			&i.IsActive,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.CreatedBy,
			&i.LastLogin,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listUsersDescending = `-- name: ListUsersDescending :many
SELECT email, is_admin, is_active, created_at, updated_at, created_by, last_login
FROM users
ORDER BY COALESCE(NULLIF($3, created_at), created_at) DESC
LIMIT $1 OFFSET $2
`

type ListUsersDescendingParams struct {
	Limit   int32       `json:"limit"`
	Offset  int32       `json:"offset"`
	Column3 interface{} `json:"column_3"`
}

type ListUsersDescendingRow struct {
	Email     string         `json:"email"`
	IsAdmin   bool           `json:"is_admin"`
	IsActive  bool           `json:"is_active"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	CreatedBy sql.NullString `json:"created_by"`
	LastLogin sql.NullTime   `json:"last_login"`
}

func (q *Queries) ListUsersDescending(ctx context.Context, arg ListUsersDescendingParams) ([]*ListUsersDescendingRow, error) {
	rows, err := q.db.Query(ctx, listUsersDescending, arg.Limit, arg.Offset, arg.Column3)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*ListUsersDescendingRow
	for rows.Next() {
		var i ListUsersDescendingRow
		if err := rows.Scan(
			&i.Email,
			&i.IsAdmin,
			&i.IsActive,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.CreatedBy,
			&i.LastLogin,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateUserAttributes = `-- name: UpdateUserAttributes :one
UPDATE users SET
                 passhash = COALESCE(NULLIF($1::varchar, ''), passhash),
                 is_active = COALESCE($2::boolean, is_active),
                 is_admin = COALESCE($3::boolean, is_admin)
WHERE email = $4
  AND (
            NULLIF($1::varchar, '') IS NOT NULL AND $1 IS DISTINCT FROM passhash OR
            $2::boolean IS NOT NULL AND $2::boolean IS DISTINCT FROM is_active OR
            $3::boolean IS NOT NULL AND $3::boolean IS DISTINCT FROM is_admin
    )
RETURNING email, is_active, is_admin, updated_at
`

type UpdateUserAttributesParams struct {
	Passhash string `json:"passhash"`
	IsActive bool   `json:"is_active"`
	IsAdmin  bool   `json:"is_admin"`
	Email    string `json:"email"`
}

type UpdateUserAttributesRow struct {
	Email     string    `json:"email"`
	IsActive  bool      `json:"is_active"`
	IsAdmin   bool      `json:"is_admin"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (q *Queries) UpdateUserAttributes(ctx context.Context, arg UpdateUserAttributesParams) (*UpdateUserAttributesRow, error) {
	row := q.db.QueryRow(ctx, updateUserAttributes,
		arg.Passhash,
		arg.IsActive,
		arg.IsAdmin,
		arg.Email,
	)
	var i UpdateUserAttributesRow
	err := row.Scan(
		&i.Email,
		&i.IsActive,
		&i.IsAdmin,
		&i.UpdatedAt,
	)
	return &i, err
}
