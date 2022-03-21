-- name: ListUsers :many
SELECT email, is_admin, is_active, created_at, updated_at, created_by, last_login
FROM users
ORDER BY COALESCE(NULLIF($3, created_at), created_at)
LIMIT $1 OFFSET $2 ;

-- name: ListUsersDescending :many
SELECT email, is_admin, is_active, created_at, updated_at, created_by, last_login
FROM users
ORDER BY COALESCE(NULLIF($3, created_at), created_at) DESC
LIMIT $1 OFFSET $2 ;

-- name: GetUserCount :one
SELECT COUNT(*) FROM users;

-- name: DeleteUser :exec
DELETE FROM users WHERE email = @email;

-- name: GetPasswordForUser :one
SELECT passhash from users where email = @email;

-- name: CreateUser :one
with u AS (INSERT INTO users (email, passhash, is_active, is_admin, created_at, updated_at, created_by)
    VALUES (@email, @passhash, @is_active, @is_admin,
            now(), now(), nullif(current_setting('request.user', true), '') )
    RETURNING *)
SELECT email, is_active, is_admin, created_at, updated_at, last_login, created_by from u;

-- name: CheckUserExists :one
SELECT count(*)::int FROM users WHERE email = @email;

-- name: GetUserByEmail :one
SELECT users.email, is_active, is_admin, created_at, updated_at, last_login, created_by
FROM users WHERE users.email = @email;

-- name: UpdateUserAttributes :one
UPDATE users SET
                 passhash = COALESCE(NULLIF(@passhash::varchar, ''), passhash),
                 is_active = COALESCE(@is_active::boolean, is_active),
                 is_admin = COALESCE(@is_admin::boolean, is_admin)
WHERE email = @email
  AND (
            NULLIF(@passhash::varchar, '') IS NOT NULL AND @passhash IS DISTINCT FROM passhash OR
            @is_active::boolean IS NOT NULL AND @is_active::boolean IS DISTINCT FROM is_active OR
            @is_admin::boolean IS NOT NULL AND @is_admin::boolean IS DISTINCT FROM is_admin
    )
RETURNING email, is_active, is_admin, updated_at;
